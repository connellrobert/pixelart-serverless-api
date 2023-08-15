package types

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ResultRequest struct {
	Record QueueRequest         `json:"record"`
	Result ImageResponseWrapper `json:"result"`
}

func (r *ResultRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: r.Record.Id,
		},
		"priority": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.Record.Priority),
		},
		"request": &types.AttributeValueMemberM{
			Value: r.Record.ToDynamoDB(),
		},
		"result": &types.AttributeValueMemberM{
			Value: r.Result.ToDynamoDB(),
		},
	}
}

func (r *ResultRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Record.Id = item["id"].(*types.AttributeValueMemberS).Value
	r.Record.Priority, _ = strconv.Atoi(item["priority"].(*types.AttributeValueMemberN).Value)
	action, err := strconv.Atoi(item["request"].(*types.AttributeValueMemberM).Value["action"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		panic(err)
	}
	r.Record.Metadata.TraceId = item["request"].(*types.AttributeValueMemberM).Value["metadata"].(*types.AttributeValueMemberM).Value["traceId"].(*types.AttributeValueMemberS).Value

	r.Record.Action = RequestAction(action)
	switch r.Record.Action {
	case GenerateImageAction:
		r.Record.CreateImage.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value)
	case EditImageAction:
		r.Record.CreateImageEdit.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value)
	case VariateImageAction:
		r.Record.CreateImageVariation.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["createImageVariation"].(*types.AttributeValueMemberM).Value)
	}
	r.Result.FromDynamoDB(item["result"].(*types.AttributeValueMemberM).Value)
}

func (r *ResultRequest) ToString() string {
	s, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(s)
}
