package main

import (
	"context"
	"fmt"

	lib "github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sashabaranov/go-openai"
)

// lambda handler
func Handler(ctx context.Context, request lib.QueueRequest) {
	var response openai.ImageResponse
	switch request.Action {
	case lib.GenerateImageAction:
		response = lib.GenerateImage(request.CreateImage)
	case lib.EditImageAction:
		response = lib.EditImage(request.EditImage)
	case lib.VariateImageAction:
		response = lib.CreateImageVariation(request.CreateImageVariation)
	}
	fmt.Println(response)
}

func main() {
	lambda.Start(Handler)
}
