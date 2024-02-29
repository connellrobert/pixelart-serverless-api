package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	aiTypes "github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
	"github.com/stretchr/testify/mock"
)

type MockSubProcess struct {
	mock.Mock
}

func (s MockSubProcess) MapParams(qr *aiTypes.QueueRequest, args aiTypes.RequestAction, params map[string]interface{}) {
	s.Called(qr, args, params)
}

func (s MockSubProcess) ToDynamoDB(obj interface {
	ToDynamoDB() map[string]types.AttributeValue
}) map[string]types.AttributeValue {
	return obj.ToDynamoDB()
}

func init() {
	subc = MockSubProcess{}
}

func TestParseSQSEvent(t *testing.T) {
	rr := aiTypes.ResultRequest{
		Record: aiTypes.QueueRequest{
			Metadata: aiTypes.CommonMetadata{
				TraceId: "1-5f9b0b9b-1",
			},
			Id:       "1234",
			Action:   aiTypes.GenerateImageAction,
			Priority: 0,
			CreateImage: aiTypes.GenerateImageRequest{
				Prompt:         "A painting of a cat",
				N:              1,
				Size:           "512x512",
				ResponseFormat: "URL",
				User:           "1234",
			},
		},
	}
	qrBytes, err := json.Marshal(rr)
	if err != nil {
		t.Fatalf("Error marshalling queue request: %s", err)
	}
	event := events.SQSMessage{
		Body: string(qrBytes),
	}
	result := ParseSQSEvent(event)
	fmt.Printf("Result: %+v\n", result)
	if result.Record.Action != aiTypes.GenerateImageAction {
		t.Fatalf("Action should be GenerateImageAction but got %d", result.Record.Action)
	}
	if result.Record.Priority != 0 {
		t.Fatalf("Priority should be 0 but got %d", result.Record.Priority)
	}
	if result.Record.CreateImage.Prompt != "A painting of a cat" {
		t.Fatalf("Prompt should be A painting of a cat but got %s", result.Record.CreateImage.Prompt)
	}
	if result.Record.CreateImage.N != 1 {
		t.Fatalf("N should be 1 but got %d", result.Record.CreateImage.N)
	}
	if result.Record.CreateImage.Size != "512x512" {
		t.Fatalf("Size should be 512x512 but got %s", result.Record.CreateImage.Size)
	}
	if result.Record.CreateImage.ResponseFormat != "URL" {
		t.Fatalf("ResponseFormat should be URL but got %s", result.Record.CreateImage.ResponseFormat)
	}
	if result.Record.CreateImage.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", result.Record.CreateImage.User)
	}
}

func TestGetAnalyticsItemInputStruct(t *testing.T) {
	tableName := "test-table"
	input := GetAnalyticsItemInputStruct("1234", tableName)
	if *input.TableName != tableName {
		t.Fatalf("TableName should be test-table but got %s", *input.TableName)
	}
	if input.Key["id"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("Key should be 1234 but got %s", input.Key["id"].(*types.AttributeValueMemberS).Value)
	}
}
