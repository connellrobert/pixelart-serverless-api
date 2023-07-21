package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib/ai"
	"github.com/aimless-it/ai-canvas/functions/lib/process"
	"github.com/aimless-it/ai-canvas/functions/lib/types"
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
	var request types.QueueRequest
	err := json.Unmarshal([]byte(queueRequest.Records[0].Body), &request)
	if err != nil {
		panic(err)
	}

	var response openai.ImageResponse
	if debug := os.Getenv("DEBUG_MODE"); debug == "true" {
		fmt.Println("DEBUG MODE IS ACTIVE")
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

		process.SendResult(request, response)
		process.SubmitXRayTraceSubSegment(request.Metadata.TraceId, "Sent result to queue")
		return
	}
	var success bool
	switch request.Action {
	case types.GenerateImageAction:
		fmt.Println("Generating image")
		response, err = ai.GenerateImage(request.CreateImage)
	case types.EditImageAction:
		fmt.Println("Editing image")
		response, err = ai.EditImage(request.CreateImageEdit)
	case types.VariateImageAction:
		fmt.Println("Varying image")
		response, err = ai.CreateImageVariation(request.CreateImageVariation)
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
	wrapped := types.ImageResponseWrapper{
		Success:  success,
		Response: response,
	}

	process.SendResult(request, wrapped)
	process.SubmitXRayTraceSubSegment(request.Metadata.TraceId, "Sent result to queue")
}

func main() {
	lambda.Start(Handler)
}
