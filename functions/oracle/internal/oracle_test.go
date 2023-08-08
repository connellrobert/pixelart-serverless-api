package internal

import (
	"encoding/json"
	"testing"

	"github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-lambda-go/events"
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

func (m mockSubprocess) GenerateImage(request types.GenerateImageRequest) (openai.ImageResponse, error) {
	args := m.Called(request)
	return args.Get(0).(openai.ImageResponse), args.Error(1)
}

func (m mockSubprocess) EditImage(request types.EditImageRequest) (openai.ImageResponse, error) {
	args := m.Called(request)
	return args.Get(0).(openai.ImageResponse), args.Error(1)
}

func (m mockSubprocess) CreateImageVariation(request types.CreateImageVariantRequest) (openai.ImageResponse, error) {
	args := m.Called(request)
	return args.Get(0).(openai.ImageResponse), args.Error(1)
}

func init() {
	subc = &mockSubprocess{}
}

func TestDeserializeSQSRequest(t *testing.T) {
	qrString, err := json.Marshal(sampleQueueRequest)
	if err != nil {
		t.Fatal(err)
	}
	qr := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(qrString),
			},
		},
	}
	result := DeserializeSQSRequest(qr)
	if result != sampleQueueRequest {
		t.Fatal("Failed to deserialize sqs request")
	}
}

func TestTestMode(t *testing.T) {
	result := TestMode()
	if result.Response.Created != sampleImageResponse.Response.Created {
		t.Fatal("Failed to return test mode response")
	}
	if result.Response.Data[0].URL != sampleImageResponse.Response.Data[0].URL {
		t.Fatal("Failed to return test mode response")
	}
	if result.Success != sampleImageResponse.Success {
		t.Fatal("Failed to return test mode response")
	}
}

func TestGenerateImage(t *testing.T) {
	subc.(*mockSubprocess).On("GenerateImage", sampleQueueRequest.CreateImage).Return(sampleImageResponse.Response, nil)
	result := GenerateImage(sampleQueueRequest.CreateImage)
	if result.Response.Created != sampleImageResponse.Response.Created {
		t.Fatal("Failed to return test mode response")
	}
	if result.Response.Data[0].URL != sampleImageResponse.Response.Data[0].URL {
		t.Fatal("Failed to return test mode response")
	}
	if result.Success != sampleImageResponse.Success {
		t.Fatal("Failed to return test mode response")
	}
}

func TestEditImage(t *testing.T) {
	subc.(*mockSubprocess).On("EditImage", sampleQueueRequest.CreateImageEdit).Return(sampleImageResponse.Response, nil)
	result := EditImage(sampleQueueRequest.CreateImageEdit)
	if result.Response.Created != sampleImageResponse.Response.Created {
		t.Fatal("Failed to return test mode response")
	}
	if result.Response.Data[0].URL != sampleImageResponse.Response.Data[0].URL {
		t.Fatal("Failed to return test mode response")
	}
	if result.Success != sampleImageResponse.Success {
		t.Fatal("Failed to return test mode response")
	}
}

func TestCreateImageVariation(t *testing.T) {
	subc.(*mockSubprocess).On("CreateImageVariation", sampleQueueRequest.CreateImageVariation).Return(sampleImageResponse.Response, nil)
	result := CreateImageVariation(sampleQueueRequest.CreateImageVariation)
	if result.Response.Created != sampleImageResponse.Response.Created {
		t.Fatal("Failed to return test mode response")
	}
	if result.Response.Data[0].URL != sampleImageResponse.Response.Data[0].URL {
		t.Fatal("Failed to return test mode response")
	}
	if result.Success != sampleImageResponse.Success {
		t.Fatal("Failed to return test mode response")
	}
}

func TestAIImageControllerGenerateImage(t *testing.T) {
	subc.(*mockSubprocess).On("GenerateImage", sampleQueueRequest.CreateImage).Return(sampleImageResponse.Response, nil)
	subc.(*mockSubprocess).On("EditImage", sampleQueueRequest.CreateImageEdit).Return(sampleImageResponse.Response, nil)
	subc.(*mockSubprocess).On("CreateImageVariation", sampleQueueRequest.CreateImageVariation).Return(sampleImageResponse.Response, nil)
	for _, action := range []types.RequestAction{types.GenerateImageAction, types.EditImageAction, types.VariateImageAction} {
		sampleQueueRequest.Action = action
		result := AIImageController(sampleQueueRequest)
		if action == types.GenerateImageAction {
			if result.Response.Created != sampleImageResponse.Response.Created {
				t.Fatal("Failed to return test mode response")
			}
			if result.Response.Data[0].URL != sampleImageResponse.Response.Data[0].URL {
				t.Fatal("Failed to return test mode response")
			}
			if result.Success != sampleImageResponse.Success {
				t.Fatal("Failed to return test mode response")
			}
		}
	}
}
