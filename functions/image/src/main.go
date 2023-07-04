package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	lambda.Start(Handler)
}

func Handler() (events.APIGatewayProxyResponse, error) {
	// Create a presigned url for s3
	svc := s3.New(s3.Options{
		Region: "us-east-1",
	})
	presign := s3.NewPresignClient(svc)
	req, er := presign.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("ai-canvas"),
		Key:    aws.String("test"),
	})
	if er != nil {
		fmt.Println("Failed to create request", er)
	}
	response, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal request", er)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
	}, nil
}
