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

func (q *QueueRequest) MapParams(action RequestAction, params interface{}) {
	switch action {
	case GenerateImageAction:
		q.CreateImage = GenerateImageRequest{
			Prompt:         params.(map[string]interface{})["prompt"].(string),
			N:              int(params.(map[string]interface{})["n"].(float64)),
			Size:           ImageSize(params.(map[string]interface{})["size"].(string)),
			ResponseFormat: ResponseFormat(params.(map[string]interface{})["responseFormat"].(string)),
			User:           params.(map[string]interface{})["user"].(string),
		}
	case EditImageAction:
		q.CreateImageEdit = EditImageRequest{
			Prompt:         params.(map[string]interface{})["prompt"].(string),
			N:              int(params.(map[string]interface{})["n"].(float64)),
			Size:           ImageSize(params.(map[string]interface{})["size"].(string)),
			ResponseFormat: ResponseFormat(params.(map[string]interface{})["responseFormat"].(string)),
			User:           params.(map[string]interface{})["user"].(string),
			Image:          params.(map[string]interface{})["image"].(string),
			Mask:           params.(map[string]interface{})["mask"].(string),
		}
	case VariateImageAction:
		q.CreateImageVariation = CreateImageVariantRequest{
			N:              int(params.(map[string]interface{})["n"].(float64)),
			Size:           ImageSize(params.(map[string]interface{})["size"].(string)),
			ResponseFormat: ResponseFormat(params.(map[string]interface{})["responseFormat"].(string)),
			User:           params.(map[string]interface{})["user"].(string),
			Image:          params.(map[string]interface{})["image"].(string),
		}
	}
}
