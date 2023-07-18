package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
	fmt.Println(queueRequest)
	err := json.Unmarshal([]byte(queueRequest.Records[0].Body), &request)
	if err != nil {
		panic(err)
	}

	var response openai.ImageResponse
	if debug := os.Getenv("DEBUG_MODE"); debug == "true" {
		fmt.Println("DEBUG MODE IS ACTIVE")
		response := lib.ImageResponseWrapper{
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

		lib.SendResult(request, response)
		lib.SubmitXRayTraceSubSegment(request.Metadata.TraceId, "Sent result to queue")
		return
	}
	var success bool
	fmt.Println(request)
	switch request.Action {
	case lib.GenerateImageAction:
		fmt.Println("Generating image")
		response, err = lib.GenerateImage(request.CreateImage)
	case lib.EditImageAction:
		fmt.Println("Editing image")
		response, err = lib.EditImage(request.CreateImageEdit)
	case lib.VariateImageAction:
		fmt.Println("Varying image")
		response, err = lib.CreateImageVariation(request.CreateImageVariation)
	default:
		fmt.Println("Invalid action")
	}
	if err != nil {
		fmt.Println(err)
		success = false
	} else {
		fmt.Println("Success!")
		success = true
	}
	fmt.Println(response)
	wrapped := lib.ImageResponseWrapper{
		Success:  success,
		Response: response,
	}

	lib.SendResult(request, wrapped)
	lib.SubmitXRayTraceSubSegment(request.Metadata.TraceId, "Sent result to queue")
}

func main() {
	lambda.Start(Handler)
}
