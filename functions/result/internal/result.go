package internal

import (
	"encoding/json"

	aiTypes "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type subProcess interface {
	ToDynamoDB(obj interface {
		ToDynamoDB() map[string]types.AttributeValue
	}) map[string]types.AttributeValue
}

type subproc struct{}

func (s subproc) ToDynamoDB(obj interface {
	ToDynamoDB() map[string]types.AttributeValue
}) map[string]types.AttributeValue {
	return obj.ToDynamoDB()
}

var subc subProcess = subproc{}

func ParseSQSEvent(message events.SQSMessage) aiTypes.ResultRequest {
	var result aiTypes.ResultRequest
	err := json.Unmarshal([]byte(message.Body), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func GetAnalyticsItemInputStruct(id, tableName string) *dynamodb.GetItemInput {

	return &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	}

}

func GetUpdateAnalyticsItemInput(analyticsItem aiTypes.AnalyticsItem, tableName string) *dynamodb.UpdateItemInput {

	return &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: analyticsItem.Id,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#A": "record",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberM{Value: subc.ToDynamoDB(&analyticsItem)},
		},
		UpdateExpression: aws.String("SET #A = :a"),
	}

}
