package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"os"
)

const PATTERN = "/v1/url-info/"

type UrlChecker interface {
	EvaluateSafety(string) (CheckerResponse, error)
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Info().Interface("request", request).Msg("Received")

	url := request.Path[len(PATTERN):]

	tableName := os.Getenv("MC_TABLE_NAME")
	if tableName == "" {
		tableName = "malicious-urls"
	}
	dataService := NewDynamoDataService(tableName)
	checker := NewMaliciousUrlChecker(dataService)

	safety, err := checker.EvaluateSafety(url)
	if err != nil {
		log.Err(err).Send()
		return events.APIGatewayProxyResponse{}, err
	}

	output, err := json.Marshal(safety)
	if err != nil {
		log.Err(err).Send()
		return events.APIGatewayProxyResponse{}, err
	}

	log.Info().Msg("Successfully processed")
	return events.APIGatewayProxyResponse{Body: string(output), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
