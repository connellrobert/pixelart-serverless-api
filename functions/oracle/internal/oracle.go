package internal

import (
	"encoding/json"
	"fmt"

	"github.com/aimless-it/ai-canvas/functions/lib/ai"
	"github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sashabaranov/go-openai"
)

type subprocess interface {
	GenerateImage(request types.GenerateImageRequest) (openai.ImageResponse, error)
	EditImage(request types.EditImageRequest) (openai.ImageResponse, error)
	CreateImageVariation(request types.CreateImageVariantRequest) (openai.ImageResponse, error)
}

func (s subproc) GenerateImage(request types.GenerateImageRequest) (openai.ImageResponse, error) {
	return ai.GenerateImage(request)
}

func (s subproc) EditImage(request types.EditImageRequest) (openai.ImageResponse, error) {
	return ai.EditImage(request)
}

func (s subproc) CreateImageVariation(request types.CreateImageVariantRequest) (openai.ImageResponse, error) {
	return ai.CreateImageVariation(request)
}

type subproc struct{}

var subc subprocess = subproc{}

func DeserializeSQSRequest(queueRequest events.SQSEvent) types.QueueRequest {
	var request types.QueueRequest
	err := json.Unmarshal([]byte(queueRequest.Records[0].Body), &request)
	if err != nil {
		panic(err)
	}
	return request
}

func TestMode() types.ImageResponseWrapper {
	fmt.Println("TEST MODE IS ACTIVE")
	response := types.ImageResponseWrapper{
		Success: true,
		Response: openai.ImageResponse{
			Created: 45454569420,
			Data: []openai.ImageResponseDataInner{
				{
					URL: "something cool",
				},
			},
		},
	}

	return response
}

func GenerateImage(request types.GenerateImageRequest) types.ImageResponseWrapper {
	response, err := subc.GenerateImage(request)
	if err != nil {
		fmt.Println(err)
		return types.ImageResponseWrapper{
			Success: false,
		}
	}
	return types.ImageResponseWrapper{
		Success:  true,
		Response: response,
	}
}

func EditImage(request types.EditImageRequest) types.ImageResponseWrapper {
	response, err := subc.EditImage(request)
	if err != nil {
		fmt.Println(err)
		return types.ImageResponseWrapper{
			Success: false,
		}
	}
	return types.ImageResponseWrapper{
		Success:  true,
		Response: response,
	}
}

func CreateImageVariation(request types.CreateImageVariantRequest) types.ImageResponseWrapper {
	response, err := subc.CreateImageVariation(request)
	if err != nil {
		fmt.Println(err)
		return types.ImageResponseWrapper{
			Success: false,
		}
	}
	return types.ImageResponseWrapper{
		Success:  true,
		Response: response,
	}
}

func AIImageController(request types.QueueRequest) types.ImageResponseWrapper {
	irw := types.ImageResponseWrapper{}
	switch request.Action {
	case types.GenerateImageAction:
		fmt.Println("Generating image")
		return GenerateImage(request.CreateImage)
	case types.EditImageAction:
		fmt.Println("Editing image")
		return EditImage(request.CreateImageEdit)
	case types.VariateImageAction:
		fmt.Println("Varying image")
		return CreateImageVariation(request.CreateImageVariation)
	default:
		fmt.Println("Invalid action")
		irw.Success = false
		return irw
	}
}
