package models

import "github.com/aws/aws-lambda-go/events"

type RespApi struct {
	Status     int
	Message    string
	CustomResp *events.APIGatewayProxyResponse
}
