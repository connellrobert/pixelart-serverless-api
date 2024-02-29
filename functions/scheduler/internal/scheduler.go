package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	aiTypes "github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
)

type subprocess interface {
	MapParams(qr *aiTypes.QueueRequest, args aiTypes.RequestAction, params map[string]interface{})
	ToDynamoDB(obj interface {
		ToDynamoDB() map[string]types.AttributeValue
	}) map[string]types.AttributeValue
}

type sub struct{}

func (s sub) MapParams(qr *aiTypes.QueueRequest, args aiTypes.RequestAction, params map[string]interface{}) {
	qr.MapParams(args, params)
}

func (s sub) ToDynamoDB(obj interface {
	ToDynamoDB() map[string]types.AttributeValue
}) map[string]types.AttributeValue {
	return obj.ToDynamoDB()
}

var subc subprocess = sub{}

func ParseApiRequest(request events.APIGatewayProxyRequest) map[string]interface{} {
	fmt.Printf("Parsing API request: %+v\n", request)
	body := make(map[string]interface{})
	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		panic(err)
	}
	return body
}

func ParseRequestAction(body map[string]interface{}) aiTypes.RequestAction {

	action, err := strconv.Atoi(fmt.Sprintf("%v", body["action"]))
	if err != nil {
		panic(err)
	}
	return aiTypes.RequestAction(action)

}

func ConvertFloatToInt(value interface{}) int {
	switch value.(type) {
	case float64:
		return int(value.(float64))
	case int:
		return value.(int)
	case string:
		i, err := strconv.Atoi(value.(string))
		if err != nil {
			panic(err)
		}
		return i
	default:
		panic("Invalid type")
	}
}

type QueueRequestArgs struct {
	Id      string
	Action  aiTypes.RequestAction
	Params  map[string]interface{}
	TraceId string
}

func ConstructQueueRequest(args QueueRequestArgs) aiTypes.QueueRequest {
	fmt.Printf("Constructing queue request: %+v\n", args)
	record := aiTypes.QueueRequest{
		Metadata: aiTypes.CommonMetadata{
			TraceId: args.TraceId,
		},
		Id:       args.Id,
		Action:   args.Action,
		Priority: 0,
	}
	subc.MapParams(&record, args.Action, args.Params)
	fmt.Printf("Constructed queue request: %+v\n", record)
	return record
}

func PutAnalyticsItem(ai aiTypes.AnalyticsItem, tableName string, dbClient *dynamodb.Client) {

	putAItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      subc.ToDynamoDB(&ai),
	}
	_, err := dbClient.PutItem(context.Background(), putAItemInput)
	if err != nil {
		panic(err)
	}
}

func ApiResponse(ai aiTypes.AnalyticsItem) (events.APIGatewayProxyResponse, error) {

	responseBody := map[string]interface{}{
		"message": fmt.Sprintf("Successfully added %v", ai.Id),
		"id":      ai.Id,
	}
	re, err := json.Marshal(responseBody)
	if err != nil {
		panic(err)
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(re),
	}
	return response, nil

}
