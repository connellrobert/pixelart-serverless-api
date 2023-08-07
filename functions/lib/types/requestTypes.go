package types

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GenerateImageRequest struct {
	Prompt         string
	N              int `default:"1" validate:"min=1,max=10"`
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
}
type EditImageRequest struct {
	Prompt         string
	N              int `default:"1" validate:"min=1,max=10"`
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
	Image          string
	Mask           string
}

type CreateImageVariantRequest struct {
	N              int `default:"1" validate:"min=1,max=10"`
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
	Image          string
}

func (r *GenerateImageRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"prompt": &types.AttributeValueMemberS{
			Value: r.Prompt,
		},
		"n": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.N),
		},
		"size": &types.AttributeValueMemberS{
			Value: string(r.Size),
		},
		"responseFormat": &types.AttributeValueMemberS{
			Value: string(r.ResponseFormat),
		},
		"user": &types.AttributeValueMemberS{
			Value: r.User,
		},
	}
}

func (r *GenerateImageRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Prompt = item["prompt"].(*types.AttributeValueMemberS).Value
	r.N, _ = strconv.Atoi(item["n"].(*types.AttributeValueMemberN).Value)
	r.Size = ImageSize(item["size"].(*types.AttributeValueMemberS).Value)
	r.ResponseFormat = ResponseFormat(item["responseFormat"].(*types.AttributeValueMemberS).Value)
	r.User = item["user"].(*types.AttributeValueMemberS).Value
}

func (r *EditImageRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"prompt": &types.AttributeValueMemberS{
			Value: r.Prompt,
		},
		"n": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.N),
		},
		"size": &types.AttributeValueMemberS{
			Value: string(r.Size),
		},
		"responseFormat": &types.AttributeValueMemberS{
			Value: string(r.ResponseFormat),
		},
		"user": &types.AttributeValueMemberS{
			Value: r.User,
		},
		"image": &types.AttributeValueMemberS{
			Value: r.Image,
		},
		"mask": &types.AttributeValueMemberS{
			Value: r.Mask,
		},
	}
}

func (r *EditImageRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Prompt = item["prompt"].(*types.AttributeValueMemberS).Value
	r.N, _ = strconv.Atoi(item["n"].(*types.AttributeValueMemberN).Value)
	r.Size = ImageSize(item["size"].(*types.AttributeValueMemberS).Value)
	r.ResponseFormat = ResponseFormat(item["responseFormat"].(*types.AttributeValueMemberS).Value)
	r.User = item["user"].(*types.AttributeValueMemberS).Value
	r.Image = item["image"].(*types.AttributeValueMemberS).Value
	r.Mask = item["mask"].(*types.AttributeValueMemberS).Value
}

func (r *CreateImageVariantRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"n": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.N),
		},
		"size": &types.AttributeValueMemberS{
			Value: string(r.Size),
		},
		"responseFormat": &types.AttributeValueMemberS{
			Value: string(r.ResponseFormat),
		},
		"user": &types.AttributeValueMemberS{
			Value: r.User,
		},
		"image": &types.AttributeValueMemberS{
			Value: r.Image,
		},
	}
}

func (r *CreateImageVariantRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.N, _ = strconv.Atoi(item["n"].(*types.AttributeValueMemberN).Value)
	r.Size = ImageSize(item["size"].(*types.AttributeValueMemberS).Value)
	r.ResponseFormat = ResponseFormat(item["responseFormat"].(*types.AttributeValueMemberS).Value)
	r.User = item["user"].(*types.AttributeValueMemberS).Value
	r.Image = item["image"].(*types.AttributeValueMemberS).Value
}
