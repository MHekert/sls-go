package shared

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func AwsEndpointResolverFactory() aws.EndpointResolverWithOptionsFunc {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		stage := os.Getenv("STAGE")

		if stage == "local" {
			return aws.Endpoint{
				URL:           "http://localhost:4566",
				PartitionID:   "aws",
				SigningRegion: "eu-central-1",
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
}
