package handler

import (
	"context"
	"errors"
	"log"
	"sls-go/src/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Key struct {
	Id string `dynamodbav:"id"`
}

func marshalKey(key *Key) map[string]types.AttributeValue {
	marshaled, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Fatal(err)
	}

	return marshaled
}

func HandlerFactory(tableName string, dynamodbClient *dynamodb.Client) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		id := event.PathParameters["id"]

		key := marshalKey(&Key{Id: id})

		dynamoResp, err := dynamodbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{Key: key, TableName: &tableName})
		if err != nil {
			var nfe *types.ResourceNotFoundException
			if errors.As(err, &nfe) {
				httpRes := shared.HTTPErrorResponse{
					StatusCode: 404,
					Message:    "not found",
				}

				return httpRes.ToAwsRes()
			}

			return shared.InternalServerError.ToAwsRes()
		}

		var item shared.Item
		err = attributevalue.UnmarshalMap(dynamoResp.Item, &item)
		if err != nil {
			log.Fatal(err)
		}

		itemJson, err := item.MarshalJson()
		if err != nil {
			log.Fatal(err)
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       itemJson,
		}, nil
	}
}
