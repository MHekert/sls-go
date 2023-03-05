package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sls-go/src/shared"
	"sls-go/src/shared/exceptions"
	"sls-go/src/shared/items"

	"github.com/aws/aws-lambda-go/events"
)

type OneGetter interface {
	GetOne(id string) (*items.Item, error)
}

func HandlerFactory(repo OneGetter) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		id := event.PathParameters["id"]

		item, err := repo.GetOne(id)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrNotFound):
				httpRes := shared.HTTPErrorResponse{
					StatusCode: 404,
					Message:    "Not Found",
				}

				return httpRes.ToAwsRes()

			default:
				fmt.Fprintln(os.Stderr, err)
				return shared.InternalServerError.ToAwsRes()
			}
		}

		itemJson, err := item.MarshalJson()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return shared.InternalServerError.ToAwsRes()
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       itemJson,
		}, nil
	}
}
