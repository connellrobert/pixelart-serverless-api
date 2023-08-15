package types

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type QueueRequest struct {
	Metadata             CommonMetadata
	Id                   string
	Action               RequestAction
	Priority             int
	CreateImage          GenerateImageRequest
	CreateImageEdit      EditImageRequest
	CreateImageVariation CreateImageVariantRequest
}

func (r *QueueRequest) ToString() string {
	json, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(json)
}
func (r *QueueRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: r.Id,
		},
		"priority": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.Priority),
		},
		"action": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(r.Action)),
		},
		"createImage": &types.AttributeValueMemberM{
			Value: r.CreateImage.ToDynamoDB(),
		},
		"createImageEdit": &types.AttributeValueMemberM{
			Value: r.CreateImageEdit.ToDynamoDB(),
		},
		"createImageVariation": &types.AttributeValueMemberM{
			Value: r.CreateImageVariation.ToDynamoDB(),
		},
		"metadata": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"traceId": &types.AttributeValueMemberS{
					Value: r.Metadata.TraceId,
				},
			},
		},
	}
}

func (r *QueueRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Id = item["id"].(*types.AttributeValueMemberS).Value
	r.Metadata.TraceId = item["metadata"].(*types.AttributeValueMemberM).Value["traceId"].(*types.AttributeValueMemberS).Value
	r.Priority, _ = strconv.Atoi(item["priority"].(*types.AttributeValueMemberN).Value)
	action, _ := strconv.Atoi(item["action"].(*types.AttributeValueMemberN).Value)
	r.Action = RequestAction(action)
	switch r.Action {
	case GenerateImageAction:
		r.CreateImage.FromDynamoDB(item["createImage"].(*types.AttributeValueMemberM).Value)
	case EditImageAction:
		r.CreateImageEdit.FromDynamoDB(item["createImageEdit"].(*types.AttributeValueMemberM).Value)
	case VariateImageAction:
		r.CreateImageVariation.FromDynamoDB(item["createImageVariation"].(*types.AttributeValueMemberM).Value)
	}
}

func (q *QueueRequest) MapParams(action RequestAction, params map[string]interface{}) {
	switch action {
	case GenerateImageAction:
		q.CreateImage = GenerateImageRequest{
			Prompt:         params["Prompt"].(string),
			N:              int(params["N"].(int)),
			Size:           ImageSize(params["Size"].(string)),
			ResponseFormat: ResponseFormat(params["ResponseFormat"].(string)),
			User:           params["User"].(string),
		}
	case EditImageAction:
		q.CreateImageEdit = EditImageRequest{
			Prompt:         params["Prompt"].(string),
			N:              int(params["N"].(int)),
			Size:           ImageSize(params["Size"].(string)),
			ResponseFormat: ResponseFormat(params["ResponseFormat"].(string)),
			User:           params["User"].(string),
			Image:          params["Image"].(string),
			Mask:           params["Mask"].(string),
		}
	case VariateImageAction:
		q.CreateImageVariation = CreateImageVariantRequest{
			N:              int(params["N"].(int)),
			Size:           ImageSize(params["Size"].(string)),
			ResponseFormat: ResponseFormat(params["ResponseFormat"].(string)),
			User:           params["User"].(string),
			Image:          params["Image"].(string),
		}
	}
}
