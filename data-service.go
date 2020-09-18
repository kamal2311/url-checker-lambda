package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)


const ITEM_NOT_FOUND = "Item not found"

type Item struct {
	Id     string
	Source string
	Score  int
}

type DynamoDataService struct {
	sess *session.Session
	svc *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDataService(tableName string) *DynamoDataService{
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return &DynamoDataService{
		sess: sess,
		svc:  svc,
		tableName: tableName,
	}
}

func (ddr *DynamoDataService) RetrieveItem(id string) (*Item, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(ddr.tableName),
	}

	result, err := ddr.svc.GetItem(input)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New(ITEM_NOT_FOUND)
	}

	dbItem := Item{}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &dbItem); err != nil {
		log.Err(err).Send()
		return nil, err
	}
	log.Info().Interface("item", dbItem).Send()
	return &dbItem, nil
}


func (ddr *DynamoDataService) InsertItem(item Item) error {
	panic("implement me")
}