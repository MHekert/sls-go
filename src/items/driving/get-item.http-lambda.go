package items

import (
	"context"
	"errors"
	"fmt"
	"os"
	service "sls-go/src/items/core/service"
	"sls-go/src/shared/common"
	"sls-go/src/shared/exceptions"

	"github.com/aws/aws-lambda-go/events"
)

func GetItemHttpLambdaHandlerFactory(getItemUseCase *service.GetItemUseCase) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		id := event.PathParameters["id"]

		item, err := getItemUseCase.Do(id)
		if err != nil {
			switch {
			case errors.Is(err, exceptions.ErrNotFound):
				httpRes := common.HTTPErrorResponse{
					StatusCode: 404,
					Message:    "Not Found",
				}

				return httpRes.ToAwsRes()

			default:
				fmt.Fprintln(os.Stderr, err)
				return common.InternalServerError.ToAwsRes()
			}
		}

		itemJson, err := item.MarshalJson()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return common.InternalServerError.ToAwsRes()
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       itemJson,
		}, nil
	}
}
