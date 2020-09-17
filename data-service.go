package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

var sess *session.Session
var svc *dynamodb.DynamoDB

func init() {

	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(sess)
}

type Item struct {
	Id     string
	Source string
	Score  int
}

func retrieveItem(id string) (*Item, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("malicious-urls"),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("Item Not Found")
	}

	dbItem := Item{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &dbItem); err != nil {
		log.Err(err).Send()
		return nil, err
	}
	log.Info().Interface("item", dbItem).Send()
	return &dbItem, nil
}
