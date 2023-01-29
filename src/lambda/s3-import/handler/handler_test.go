package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

type Context struct {
	context.Context
}

func TestSomething(t *testing.T) {

	event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventSource:  "aws:s3",
				EventVersion: "2.0",
				AWSRegion:    "eu-central-1",
				EventName:    "ObjectCreated:Put",
				EventTime:    time.Now(),
				PrincipalID: events.S3UserIdentity{
					PrincipalID: "EXAMPLE",
				},
				S3: events.S3Entity{
					SchemaVersion:   "1.0",
					ConfigurationID: "testConfigRule",
					Bucket: events.S3Bucket{
						Name: "sls-go-local-123456789-files-s3-bucket",
						OwnerIdentity: events.S3UserIdentity{
							PrincipalID: "EXAMPLE",
						},
						Arn: "arn:aws:s3:::sls-go-local-123456789-files-s3-bucket",
					},
					Object: events.S3Object{
						Key:           "uploads/example.csv",
						Size:          1024,
						URLDecodedKey: "uploads/example.csv",
						VersionID:     "1",
						ETag:          "0123456789abcdef0123456789abcdef",
						Sequencer:     "0A1B2C3D4E5F678901",
					},
				},
				RequestParameters: events.S3RequestParameters{
					SourceIPAddress: "127.0.0.1",
				},
				ResponseElements: make(map[string]string)},
		},
	}

	fmt.Println(event)
	// Handler(Context{}, event)

	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")

}
