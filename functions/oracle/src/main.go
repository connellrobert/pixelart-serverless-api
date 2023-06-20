package main

import (
	"context"
	"fmt"

	lib "github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/lambda"
	openai "github.com/sashabaranov/go-openai"
)

// lambda handler
func Handler(ctx context.Context, request lib.QueueRequest) {
	var response openai.ImageResponse
	switch request.Action {
	case lib.GenerateImageAction:
		response = lib.GenerateImage(request.CreateImage)
	case lib.EditImageAction:
		response = lib.EditImage(request.CreateImageEdit)
	case lib.VariateImageAction:
		response = lib.CreateImageVariation(request.CreateImageVariation)
	}
	fmt.Println(response)
	wrapped := lib.ImageResponseWrapper{
		Response: response,
	}
	lib.SendResult(request, wrapped)
}

func main() {
	lambda.Start(Handler)
}
