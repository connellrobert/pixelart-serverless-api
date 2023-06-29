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
	err := json.Unmarshal([]byte(queueRequest.Records[0].Body), &request)
	if err != nil {
		panic(err)
	}

	var response openai.ImageResponse
	if debug := os.Getenv("DEBUG_MODE"); debug == "true" {
		response := lib.ImageResponseWrapper{
			Success: true,
			Response: openai.ImageResponse{
				Created: 45454569420,
				Data: []openai.ImageResponseDataInner{
					openai.ImageResponseDataInner{
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
	switch request.Action {
	case lib.GenerateImageAction:
		response, err = lib.GenerateImage(request.CreateImage)
	case lib.EditImageAction:
		response, err = lib.EditImage(request.CreateImageEdit)
	case lib.VariateImageAction:
		response, err = lib.CreateImageVariation(request.CreateImageVariation)
	}
	if err != nil {
		success = false
	} else {
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
