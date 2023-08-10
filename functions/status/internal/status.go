package internal

import (
	"context"
	"encoding/json"
	"fmt"

	aiTypes "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetAnalyticsItem(id, tableName string, client *dynamodb.Client) map[string]types.AttributeValue {
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	}
	record, err := client.GetItem(context.Background(), getItemInput)
	if err != nil {
		panic(err)
	}
	return record.Item
}

func GetAnalyticsItemAttemptsUrls(attempts map[string]aiTypes.ImageResponseWrapper) []string {
	urls := []string{}
	for _, attempt := range attempts {
		for _, url := range attempt.Response.Data {
			urls = append(urls, url.URL)
		}
	}
	return urls
}

func CreateResponse(urls []string, ai aiTypes.AnalyticsItem) (events.APIGatewayProxyResponse, error) {
	if len(urls) == 0 {
		if len(ai.Attempts) >= 3 {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "{\"message\": \"No successful attempts\"}",
			}, nil
		}
		// return empty message
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
		}, nil
	}
	jUrl, err := json.Marshal(urls)
	if err != nil {
		panic(err)
	}
	// return analytics item
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("{\"urls\": \" %s \"}", jUrl),
	}, nil

}
