package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Handler(ctx context.Context, s3Event events.S3Event) {
	fmt.Println("inside")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	for _, record := range s3Event.Records {
		s3data := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3data.Bucket.Name, s3data.Object.Key)

		fd, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s3data.Bucket.Name),
			Key:    aws.String(s3data.Object.Key),
		})
		if err != nil {
			panic(err)
		}
		defer fd.Body.Close()

		csvReader := csv.NewReader(fd.Body)
		for {
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			fmt.Println(row)
		}

	}

}
