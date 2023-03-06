package common

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type HTTPErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

var InternalServerError = HTTPErrorResponse{
	StatusCode: 500,
	Message:    "Internal Server Error",
}

func (res *HTTPErrorResponse) MarshalJson() (string, error) {
	jsonString, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}

func (res *HTTPErrorResponse) ToAwsRes() (events.APIGatewayProxyResponse, error) {
	httpResStr, err := res.MarshalJson()
	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: res.StatusCode,
		Body:       httpResStr,
	}, nil
}
