package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
)

var (
	sampleQueueRequest = types.QueueRequest{
		Metadata: types.CommonMetadata{
			TraceId: "1-5f9b0b9b-1",
		},
		CreateImage: types.GenerateImageRequest{
			Prompt:         "A painting of a cat",
			N:              1,
			Size:           "512x512",
			ResponseFormat: "URL",
			User:           "1234",
		},
		Id:       "1234",
		Action:   types.GenerateImageAction,
		Priority: 1,
		CreateImageEdit: types.EditImageRequest{
			Prompt:         "A painting of a cat",
			N:              1,
			Size:           "512x512",
			ResponseFormat: "URL",
			User:           "1234",
			Image:          "test",
			Mask:           "test",
		},
		CreateImageVariation: types.CreateImageVariantRequest{
			N:              1,
			Size:           "512x512",
			ResponseFormat: "URL",
			User:           "1234",
			Image:          "test",
		},
	}

	sampleImageResponse = types.ImageResponseWrapper{
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
)

type mockSubprocess struct {
	mock.Mock
}

func (m mockSubprocess) DeserializeSQSRequest(queueRequest events.SQSEvent) types.QueueRequest {
	args := m.Called(queueRequest)
	return args.Get(0).(types.QueueRequest)
}

func (m mockSubprocess) TestMode() types.ImageResponseWrapper {
	args := m.Called()
	return args.Get(0).(types.ImageResponseWrapper)
}

func (m mockSubprocess) AIImageController(request types.QueueRequest) types.ImageResponseWrapper {
	args := m.Called(request)
	return args.Get(0).(types.ImageResponseWrapper)
}

func (m mockSubprocess) SendResult(request types.QueueRequest, wrapped types.ImageResponseWrapper) {
	m.Called(request, wrapped)
}

func (m mockSubprocess) SubmitXRayTraceSubSegment(traceId string, name string) {
	m.Called(traceId, name)
}

func (m mockSubprocess) SaveImage(base64 string, fileName string) {
	m.Called(base64, fileName)
}

func (m mockSubprocess) GetPresignedUrl(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func init() {
	subprocessor = &mockSubprocess{}
}

func TestHandler(t *testing.T) {
	os.Setenv("MODE", "test")
	queueRequest := events.SQSEvent{
		Records: []events.SQSMessage{
			{},
		},
	}
	subprocessor.(*mockSubprocess).On("DeserializeSQSRequest", queueRequest).Return(sampleQueueRequest)
	subprocessor.(*mockSubprocess).On("TestMode").Return(sampleImageResponse)
	subprocessor.(*mockSubprocess).On("SendResult", sampleQueueRequest, sampleImageResponse)
	subprocessor.(*mockSubprocess).On("SubmitXRayTraceSubSegment", sampleQueueRequest.Metadata.TraceId, "Sent result to queue")
	Handler(context.Background(), queueRequest)
}

func TestHandlerProd(t *testing.T) {
	os.Setenv("MODE", "")
	queueRequest := events.SQSEvent{
		Records: []events.SQSMessage{
			{},
		},
	}
	subprocessor.(*mockSubprocess).On("DeserializeSQSRequest", queueRequest).Return(sampleQueueRequest)
	subprocessor.(*mockSubprocess).On("AIImageController", sampleQueueRequest).Return(sampleImageResponse)
	subprocessor.(*mockSubprocess).On("SendResult", sampleQueueRequest, sampleImageResponse)
	subprocessor.(*mockSubprocess).On("SubmitXRayTraceSubSegment", sampleQueueRequest.Metadata.TraceId, "Sent result to queue")
	Handler(nil, queueRequest)
}

func TestHandlerProdError(t *testing.T) {
	os.Setenv("MODE", "")
	queueRequest := events.SQSEvent{
		Records: []events.SQSMessage{
			{},
		},
	}
	subprocessor.(*mockSubprocess).On("DeserializeSQSRequest", queueRequest).Return(sampleQueueRequest)
	subprocessor.(*mockSubprocess).On("AIImageController", sampleQueueRequest).Return(types.ImageResponseWrapper{
		Success: false,
	})
	subprocessor.(*mockSubprocess).On("SendResult", sampleQueueRequest, types.ImageResponseWrapper{
		Success: false,
	})
	subprocessor.(*mockSubprocess).On("SubmitXRayTraceSubSegment", sampleQueueRequest.Metadata.TraceId, "Sent result to queue")
	Handler(nil, queueRequest)
}
