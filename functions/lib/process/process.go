package process

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

	. "github.com/aimless-it/ai-canvas/functions/lib/types"
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

func GetAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(Region()),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Region() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	return region
}

func SendResult(record QueueRequest, response ImageResponseWrapper) {
	// Create a Lambda client
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
	sqsClient := sqs.NewFromConfig(GetAWSConfig())
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
	client := dynamodb.NewFromConfig(GetAWSConfig())
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
	_, err := client.PutItem(context.Background(), putItemInput)
	if err != nil {
		panic(err)
	}
	return "success"
}

func SendRetrySignalV2(record QueueRequest) {

	j, err := json.Marshal(record)
	if err != nil {
		panic(err)
	}

	queueUrl := os.Getenv("QUEUE_URL")
	sqsClient := sqs.NewFromConfig(GetAWSConfig())

	// send item to queue
	sendMessageInput := &sqs.SendMessageInput{
		MessageBody:            aws.String(string(j)),
		QueueUrl:               aws.String(queueUrl),
		MessageGroupId:         aws.String("1"),
		MessageDeduplicationId: aws.String(record.Id),
	}
	_, err = sqsClient.SendMessage(context.Background(), sendMessageInput)
	if err != nil {
		panic(err)
	}
}

func SetAlarmState(name, status string) {
	cloudwatchClient := cloudwatch.NewFromConfig(GetAWSConfig())
	setAlarmStateInput := &cloudwatch.SetAlarmStateInput{
		AlarmName:   aws.String(name),
		StateValue:  cTypes.StateValue(status),
		StateReason: aws.String("Cause I be balling like an OG in the club"),
	}
	_, err := cloudwatchClient.SetAlarmState(context.Background(), setAlarmStateInput)
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
	id := strings.Replace(uuid.New().String(), "-", "", -1)[:16]
	xrayClient := xray.NewFromConfig(GetAWSConfig())
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	document := map[string]interface{}{
		"trace_id":   parentSegmentId,
		"id":         id,
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

func SendRequestToQueue(record QueueRequest) {
	j, err := json.Marshal(record)
	if err != nil {
		panic(err)
	}

	queueUrl := os.Getenv("QUEUE_URL")
	sqsClient := sqs.NewFromConfig(GetAWSConfig())

	// send item to queue
	sendMessageInput := &sqs.SendMessageInput{
		MessageBody:            aws.String(string(j)),
		QueueUrl:               aws.String(queueUrl),
		MessageGroupId:         aws.String("1"),
		MessageDeduplicationId: aws.String(record.Id),
	}
	_, err = sqsClient.SendMessage(context.Background(), sendMessageInput)
	if err != nil {
		panic(err)
	}
}

func StoreAnalyticsItem(ai AnalyticsItem) {
	client := dynamodb.NewFromConfig(GetAWSConfig())

	analyticsTable := os.Getenv("ANALYTICS_TABLE_NAME")
	putAItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(analyticsTable),
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: ai.Id,
			},
			"Record": &types.AttributeValueMemberM{
				Value: ai.ToDynamoDB(),
			},
		},
	}
	_, err := client.PutItem(context.Background(), putAItemInput)
	if err != nil {
		panic(err)
	}
}
