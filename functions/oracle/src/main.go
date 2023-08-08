package main

import (
	"context"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib/process"
	"github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aimless-it/ai-canvas/functions/oracle/internal"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type subprocess interface {
	DeserializeSQSRequest(queueRequest events.SQSEvent) types.QueueRequest
	TestMode() types.ImageResponseWrapper
	AIImageController(request types.QueueRequest) types.ImageResponseWrapper
	SendResult(request types.QueueRequest, wrapped types.ImageResponseWrapper)
	SubmitXRayTraceSubSegment(traceId string, name string)
}

type subproc struct{}

func (s subproc) DeserializeSQSRequest(queueRequest events.SQSEvent) types.QueueRequest {
	return internal.DeserializeSQSRequest(queueRequest)
}

func (s subproc) TestMode() types.ImageResponseWrapper {
	return internal.TestMode()
}

func (s subproc) AIImageController(request types.QueueRequest) types.ImageResponseWrapper {
	return internal.AIImageController(request)
}

func (s subproc) SendResult(request types.QueueRequest, wrapped types.ImageResponseWrapper) {
	process.SendResult(request, wrapped)
}

func (s subproc) SubmitXRayTraceSubSegment(traceId string, name string) {
	process.SubmitXRayTraceSubSegment(traceId, name)
}

var subprocessor subprocess = subproc{}

// List of environment variables:
// OPENAI_API_KEY
// RESULT_FUNCTION_ARN - the ARN of the result function (not used in this file)
// lambda handler
// TODO: Retrieve images from s3 prior to calling openai requests
func Handler(ctx context.Context, queueRequest events.SQSEvent) {
	request := subprocessor.DeserializeSQSRequest(queueRequest)
	var wrapped types.ImageResponseWrapper
	if debug := os.Getenv("MODE"); debug == "test" {
		wrapped = subprocessor.TestMode()
	} else {
		wrapped = subprocessor.AIImageController(request)
	}
	subprocessor.SendResult(request, wrapped)
	subprocessor.SubmitXRayTraceSubSegment(request.Metadata.TraceId, "Sent result to queue")
}

func main() {
	lambda.Start(Handler)
}
