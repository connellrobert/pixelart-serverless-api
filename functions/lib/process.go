package lib
// environment variables required for this file:
// RESULT_FUNCTION_ARN - the ARN of the result function
// GENERATE_IMAGE_TABLE_NAME - the name of the generate image table
// EDIT_IMAGE_TABLE_NAME - the name of the edit image table
// VARIATE_IMAGE_TABLE_NAME - the name of the variate image table

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func SendResult(record interface{}, response interface{}) string {
	// Create a Lambda client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	// invoke lambda
	lambdaArn := os.Getenv("RESULT_FUNCTION_ARN")
	svc := lambda.NewFromConfig(cfg)
	tmp := ResultRequest{
		Record:   record,
		Response: response,
	}
	req, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	input := &lambda.InvokeInput{
		FunctionName: aws.String(lambdaArn),
		Payload:      []byte(req),
	}
	result, err := svc.Invoke(context.Background(), input)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result.Payload))
	return string(result.Payload)
}

func SendRetrySignal(record QueueRequest) string {
	record.Priority = 1
	var tableName string;
	switch record.Action {
		case GenerateImageAction:
			tableName = os.Getenv("GENERATE_IMAGE_TABLE_NAME")
		case EditImageAction:
			tableName = os.Getenv("EDIT_IMAGE_TABLE_NAME")
		case VariateImageAction:
			tableName = os.Getenv("VARIATE_IMAGE_TABLE_NAME")
	
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				id: record.Id,
			},
			"Priority": &types.AttributeValueMemberN{
				id: record.Priority,
			},
		},
	}
	_, err = client.PutItem(context.Background(), putItemInput)
	if err != nil {
		panic(err)
	}
	return "success"
}
