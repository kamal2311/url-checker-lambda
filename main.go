package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const PATTERN = "/v1/url-info/"

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	url := request.Path[len(PATTERN):]

	checker := NewMaliciousUrlChecker()
	safety := checker.EvaluateSafety(url)

	output, err := json.Marshal(safety)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: string(output), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
