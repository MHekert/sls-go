package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	handler "sls-go/src/lambda/s3-import/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
