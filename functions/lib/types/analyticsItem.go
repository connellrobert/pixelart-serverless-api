package types

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AnalyticsItem struct {
	Success  bool                            `json:"success"`
	Id       string                          `json:"id"`
	Record   QueueRequest                    `json:"record"`
	Attempts map[string]ImageResponseWrapper `json:"attempts"`
}

// Create dynamodb mappings for AnalyticsItem
func (r *AnalyticsItem) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: r.Id,
		},
		"record": &types.AttributeValueMemberM{
			Value: r.Record.ToDynamoDB(),
		},
		"attempts": &types.AttributeValueMemberM{
			Value: r.AttemptsToDynamoDB(),
		},
		"success": &types.AttributeValueMemberBOOL{
			Value: r.Success,
		},
	}
}

func (r *AnalyticsItem) FromDynamoDB(item map[string]types.AttributeValue) {
	record := item["record"].(*types.AttributeValueMemberM).Value
	r.Id = item["id"].(*types.AttributeValueMemberS).Value
	r.Success = item["success"].(*types.AttributeValueMemberBOOL).Value
	// record := request["request"].(*types.AttributeValueMemberM).Value
	action, err := strconv.Atoi(record["action"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		panic(err)
	}
	r.Record.Action = RequestAction(action)
	switch r.Record.Action {
	case GenerateImageAction:
		r.Record.CreateImage.FromDynamoDB(record["createImage"].(*types.AttributeValueMemberM).Value)
	case EditImageAction:
		r.Record.CreateImageEdit.FromDynamoDB(record["createImageEdit"].(*types.AttributeValueMemberM).Value)
	case VariateImageAction:
		r.Record.CreateImageVariation.FromDynamoDB(record["createImageVariation"].(*types.AttributeValueMemberM).Value)
	}
	r.Record.Metadata = CommonMetadata{
		TraceId: record["metadata"].(*types.AttributeValueMemberM).Value["traceId"].(*types.AttributeValueMemberS).Value,
	}
	r.Record.Id = record["id"].(*types.AttributeValueMemberS).Value
	r.Record.Priority, _ = strconv.Atoi(record["priority"].(*types.AttributeValueMemberN).Value)

	r.AttemptsFromDynamoDB(item["attempts"].(*types.AttributeValueMemberM).Value)
}

func (r *AnalyticsItem) AttemptsToDynamoDB() map[string]types.AttributeValue {
	attempts := make(map[string]types.AttributeValue)
	for k, v := range r.Attempts {
		attempts[k] = &types.AttributeValueMemberM{
			Value: v.ToDynamoDB(),
		}
		attempts[k].(*types.AttributeValueMemberM).Value["success"] = &types.AttributeValueMemberBOOL{
			Value: v.Success,
		}
	}
	return attempts
}

func (r *AnalyticsItem) AttemptsFromDynamoDB(item map[string]types.AttributeValue) {
	r.Attempts = make(map[string]ImageResponseWrapper)
	for k, v := range item {
		var irw ImageResponseWrapper
		irw.FromDynamoDB(v.(*types.AttributeValueMemberM).Value)
		r.Attempts[k] = irw
	}
}
