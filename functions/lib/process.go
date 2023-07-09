package lib

// environment variables required for this file:
// RESULT_FUNCTION_ARN - the ARN of the result function
// GENERATE_IMAGE_TABLE_NAME - the name of the generate image table
// EDIT_IMAGE_TABLE_NAME - the name of the edit image table
// VARIATE_IMAGE_TABLE_NAME - the name of the variate image table

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/google/uuid"
)

func SendResult(record QueueRequest, response ImageResponseWrapper) {
	// Create a Lambda client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	// invoke lambda
	tmp := ResultRequest{
		Record: record,
		Result: response,
	}
	req, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}
	// Send req to sqs queue
	queue := os.Getenv("RESULT_QUEUE_URL")
	sqsClient := sqs.NewFromConfig(cfg)
	// send item to queue
	messageInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(req)),
		QueueUrl:    aws.String(queue),
	}
	_, err = sqsClient.SendMessage(context.Background(), messageInput)
	if err != nil {
		panic(err)
	}

}

func SendRetrySignal(record QueueRequest) string {
	record.Priority = 1
	var tableName string
	switch record.Action {
	case GenerateImageAction:
		tableName = os.Getenv("GENERATE_IMAGE_TABLE_NAME")
	case EditImageAction:
		tableName = os.Getenv("EDIT_IMAGE_TABLE_NAME")
	case VariateImageAction:
		tableName = os.Getenv("VARIATE_IMAGE_TABLE_NAME")
	}
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
			"id": &types.AttributeValueMemberS{
				Value: record.Id,
			},
			"priority": &types.AttributeValueMemberN{
				Value: strconv.Itoa(record.Priority),
			},
		},
	}
	_, err = client.PutItem(context.Background(), putItemInput)
	if err != nil {
		panic(err)
	}
	return "success"
}

func SetAlarmState(name, status string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	setAlarmStateInput := &cloudwatch.SetAlarmStateInput{
		AlarmName:   aws.String(name),
		StateValue:  cTypes.StateValue(status),
		StateReason: aws.String("Cause I be balling like an OG in the club"),
	}
	_, err = cloudwatchClient.SetAlarmState(context.Background(), setAlarmStateInput)
	if err != nil {
		panic(err)
	}

}

var ActionToTableEnvMapping = map[RequestAction]string{
	GenerateImageAction: "GI_TABLE_NAME",
	EditImageAction:     "EI_TABLE_NAME",
	VariateImageAction:  "VI_TABLE_NAME",
}

var ActionToAlarmMapping = map[RequestAction]string{
	GenerateImageAction: "giLowDynamoDBCountAlarm",
	EditImageAction:     "eiLowDynamoDBCountAlarm",
	VariateImageAction:  "viLowDynamoDBCountAlarm",
}

func SubmitXRayTraceSubSegment(parentSegmentId, name string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	id := strings.Replace(uuid.New().String(), "-", "", -1)[:16]
	xrayClient := xray.NewFromConfig(cfg)
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	document := map[string]interface{}{
		"trace_id":   id,
		"id":         parentSegmentId,
		"name":       name,
		"start_time": startTime,
		"end_time":   startTime + 10000,
	}
	documentString, err := json.Marshal(document)
	if err != nil {
		panic(err)
	}

	submitSubSegmentInput := &xray.PutTraceSegmentsInput{
		TraceSegmentDocuments: []string{
			string(documentString),
		},
	}
	_, err = xrayClient.PutTraceSegments(context.Background(), submitSubSegmentInput)
	if err != nil {
		panic(err)
	}
}
