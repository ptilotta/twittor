package models

import "github.com/aws/aws-lambda-go/events"

type RespApi struct {
	Status     int `default:"400"`
	Message    string
	CustomResp *events.APIGatewayProxyResponse
}
