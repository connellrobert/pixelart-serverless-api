package types

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	openai "github.com/sashabaranov/go-openai"
)

type ImageResponseWrapper struct {
	Response openai.ImageResponse `json:"response"`
	Success  bool                 `json:"success"`
}

// Create dynamodb AttributeValue List for openai response
func (r *ImageResponseWrapper) MapDataInnerToDynamoDB() []types.AttributeValue {
	data := make([]types.AttributeValue, 0)
	for _, d := range r.Response.Data {
		data = append(data, &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"url": &types.AttributeValueMemberS{
					Value: d.URL,
				},
				"b64": &types.AttributeValueMemberS{
					Value: d.B64JSON,
				},
			},
		})
	}
	return data
}

func (r *ImageResponseWrapper) MapDynamoDBToDataInner(dbList []types.AttributeValue) []openai.ImageResponseDataInner {
	var dataInner []openai.ImageResponseDataInner
	for _, d := range dbList {
		var url, b64 string
		if d.(*types.AttributeValueMemberM).Value["url"] == nil {
			url = ""
		} else {
			url = d.(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value
		}
		if d.(*types.AttributeValueMemberM).Value["b64"] == nil {
			b64 = ""
		} else {
			b64 = d.(*types.AttributeValueMemberM).Value["b64"].(*types.AttributeValueMemberS).Value
		}
		dataInner = append(dataInner, openai.ImageResponseDataInner{
			URL:     url,
			B64JSON: b64,
		})
	}
	return dataInner
}

func (r *ImageResponseWrapper) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"success": &types.AttributeValueMemberBOOL{
			Value: r.Success,
		},
		"response": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"created": &types.AttributeValueMemberN{
					Value: strconv.Itoa(int(r.Response.Created)),
				},
				"data": &types.AttributeValueMemberL{
					Value: r.MapDataInnerToDynamoDB(),
				},
			},
		},
	}
}

func (r *ImageResponseWrapper) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Response.Data = make([]openai.ImageResponseDataInner, 0)
	r.Success = item["success"].(*types.AttributeValueMemberBOOL).Value
	c, _ := strconv.Atoi(item["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value)
	r.Response.Created = int64(c)
	r.Response.Data = r.MapDynamoDBToDataInner(item["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value)
}
