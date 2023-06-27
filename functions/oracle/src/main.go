package main

import (
	"context"
	"encoding/json"
	"fmt"

	lib "github.com/aimless-it/ai-canvas/functions/lib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	openai "github.com/sashabaranov/go-openai"
)

// List of environment variables:
// OPENAI_API_KEY
// RESULT_FUNCTION_ARN - the ARN of the result function (not used in this file)
// lambda handler
// TODO: Retrieve images from s3 prior to calling openai requests
func Handler(ctx context.Context, queueRequest events.SQSEvent) {
	var request lib.QueueRequest
	err := json.Unmarshal([]byte(queueRequest.Records[0].Body), &request)
	if err != nil {
		panic(err)
	}

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
