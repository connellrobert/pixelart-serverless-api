package main

import (
	"context"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib/process"
	aiTypes "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// List of environment variables:
// ANALYTICS_TABLE_NAME

func main() {
	lambda.Start(Handler)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get id from request path
	id := req.PathParameters["id"]
	tableName := os.Getenv("ANALYTICS_TABLE_NAME")
	// check analytics db for id
	cfg := process.GetAWSConfig()
	client := dynamodb.NewFromConfig(cfg)
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
	analyticsItem := aiTypes.AnalyticsItem{}
	analyticsItem.FromDynamoDB(record.Item)
	process.SubmitXRayTraceSubSegment(analyticsItem.Record.Metadata.TraceId, "Retrieved analytics item from db")
	var url string
	for _, attempt := range analyticsItem.Attempts {
		if attempt.Success {
			url = attempt.Response.Data[0].URL
			break
		}
	}
	// check if data is empty
	if len(url) == 0 {
		if len(analyticsItem.Attempts) >= 3 {
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
	// return analytics item
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{\"url\": \"" + url + "\"}",
	}, nil

}
