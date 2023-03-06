package items

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type EmailType string

type Item struct {
	Id        string    `dynamodbav:"id" csv:"id" json:"id"`
	FirstName string    `dynamodbav:"firstName" csv:"firstname" json:"firstname"`
	LastName  string    `dynamodbav:"lastName" csv:"lastname" json:"lastname"`
	Email     EmailType `dynamodbav:"email" csv:"email" json:"email"`
	Value     int       `dynamodbav:"value" csv:"value" json:"value"`
}

func (item *Item) MarshalJson() (string, error) {
	jsonString, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}

func (item *Item) MarshalDynamodbav() (map[string]types.AttributeValue, error) {
	marshaledItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	return marshaledItem, nil
}
