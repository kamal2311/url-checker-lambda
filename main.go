package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

const PATTERN = "/v1/url-info/"

type UrlService interface {
	Check(string) (CheckerResponse, error)
	Save(string) error
}

func main() {
	lambda.Start(reqRouter)
}

func reqRouter(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info().Interface("request", request).Msg("Received")
	switch request.HTTPMethod {
	case "GET":
		return handleGet(ctx, request)
	case "POST":
		return handlePost(ctx, request)
	default:
		return handleClientError()
	}
}

func handleGet(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	path := request.Path[len(PATTERN):]
	log.Info().Str("urlId", path).Send()

	tableName := os.Getenv("MC_TABLE_NAME")
	if tableName == "" {
		tableName = "malicious-urls"
	}
	dataService := NewDynamoDataService(tableName)
	urlService := NewMaliciousUrlService(dataService)

	safety, err := urlService.Check(path)
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
	return events.APIGatewayProxyResponse{Body: string(output), StatusCode: http.StatusOK}, nil

}

func handlePost(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	tableName := os.Getenv("MC_TABLE_NAME")
	if tableName == "" {
		tableName = "malicious-urls"
	}

	dataService := NewDynamoDataService(tableName)
	urlService := NewMaliciousUrlService(dataService)

	item, err := validateAndDecode(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, nil
	}

	err = urlService.Save(item)
	if err != nil {
		log.Err(err).Send()
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(request.Body), StatusCode: http.StatusCreated}, nil
}



func handleClientError() (events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Allow"] = "GET, POST"
	return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed, Headers: headers}, nil
}
