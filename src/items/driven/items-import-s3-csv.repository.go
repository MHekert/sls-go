package items

import (
	"context"
	"encoding/csv"
	"io"
	items "sls-go/src/items/core"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jszwec/csvutil"
)

type ItemsImportS3CSVRepository struct {
	client     *s3.Client
	bucketName string
}

func NewItemsImportS3CSVRepository(s3 *s3.Client, bucketName string) *ItemsImportS3CSVRepository {
	return &ItemsImportS3CSVRepository{
		client:     s3,
		bucketName: bucketName,
	}
}

func (repo *ItemsImportS3CSVRepository) unmarshalEmail(data []byte, field *items.EmailType) error {
	email := strings.ToLower(string(data))
	*field = items.EmailType(email)
	return nil
}

func (repo *ItemsImportS3CSVRepository) GetImportItemsChannel(key string, importChannel chan<- items.Item) error {
	fd, err := repo.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(repo.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer fd.Body.Close()

	csvReader := csv.NewReader(fd.Body)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return err
	}
	dec.Register(repo.unmarshalEmail)

	for {
		var item items.Item
		err := dec.Decode(&item)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		importChannel <- item
	}

	close(importChannel)
	return nil
}
