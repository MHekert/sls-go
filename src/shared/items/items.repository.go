package items

import (
	"context"
	"log"
	"sls-go/src/shared/exceptions"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const MaxBatchSize = 25

type ItemsDynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

type itemKey struct {
	Id string `dynamodbav:"id"`
}

func NewItemsDynamoDBRepository(dynamodb *dynamodb.Client, tableName string) *ItemsDynamoDBRepository {
	return &ItemsDynamoDBRepository{
		client:    dynamodb,
		tableName: tableName,
	}
}

func (repo *ItemsDynamoDBRepository) PersistBatch(items [](*Item)) error {
	batchWriteInput := dynamodb.BatchWriteItemInput{}
	batchWriteInput.RequestItems = make(map[string][]types.WriteRequest)
	batchItems := make([]types.WriteRequest, 0, MaxBatchSize)

	for i := range items {
		marshaledItem, err := items[i].MarshalDynamodbav()
		if err != nil {
			return err

		}

		putRequest := types.PutRequest{
			Item: marshaledItem,
		}
		batchItems = append(batchItems, types.WriteRequest{
			PutRequest: &putRequest,
		})
	}
	batchWriteInput.RequestItems[repo.tableName] = batchItems

	_, err := repo.client.BatchWriteItem(context.TODO(), &batchWriteInput)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ItemsDynamoDBRepository) marshalKey(key *itemKey) map[string]types.AttributeValue {
	marshaled, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Fatal(err)
	}

	return marshaled
}

func (repo *ItemsDynamoDBRepository) GetOne(id string) (*Item, error) {
	key := repo.marshalKey(&itemKey{Id: id})

	dynamoResp, err := repo.client.GetItem(context.TODO(), &dynamodb.GetItemInput{Key: key, TableName: &repo.tableName})
	if err != nil {
		return nil, err
	}

	if len(dynamoResp.Item) == 0 {
		return nil, exceptions.ErrNotFound
	}

	var item Item
	err = attributevalue.UnmarshalMap(dynamoResp.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
