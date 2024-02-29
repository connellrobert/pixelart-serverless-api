package internal

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	aiTypes "github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
	"github.com/sashabaranov/go-openai"
)

func TestConstructQueueRequest(t *testing.T) {
	args := QueueRequestArgs{
		Id:     "1234",
		Action: aiTypes.GenerateImageAction,
		Params: map[string]interface{}{
			"Prompt":         "A painting of a cat",
			"N":              1,
			"Size":           "512x512",
			"ResponseFormat": "URL",
			"User":           "1234",
		},
		TraceId: "1-5f9b0b9b-1",
	}
	queueRequest := ConstructQueueRequest(args)
	if queueRequest.Id != "1234" {
		t.Fatalf("Id should be 1234 but got %s", queueRequest.Id)
	}
	if queueRequest.Action != aiTypes.GenerateImageAction {
		t.Fatalf("Action should be GenerateImageAction but got %d", queueRequest.Action)
	}
	if queueRequest.Priority != 0 {
		t.Fatalf("Priority should be 1 but got %d", queueRequest.Priority)
	}
	if queueRequest.CreateImage.Prompt != "A painting of a cat" {
		t.Fatalf("Prompt should be A painting of a cat but got %s", queueRequest.CreateImage.Prompt)
	}
	if queueRequest.CreateImage.N != 1 {
		t.Fatalf("N should be 1 but got %d", queueRequest.CreateImage.N)
	}
	if queueRequest.CreateImage.Size != "512x512" {
		t.Fatalf("Size should be 512x512 but got %s", queueRequest.CreateImage.Size)
	}
	if queueRequest.CreateImage.ResponseFormat != "URL" {
		t.Fatalf("ResponseFormat should be URL but got %s", queueRequest.CreateImage.ResponseFormat)
	}
	if queueRequest.CreateImage.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", queueRequest.CreateImage.User)
	}

}

func TestParseApiRequest(t *testing.T) {
	reqBody := map[string]interface{}{
		"action": 0,
		"params": map[string]interface{}{
			"image":          "https://aimless.ai/images/ai-canvas-logo.png",
			"size":           "512x512",
			"prompt":         "something simple",
			"n":              1,
			"responseFormat": "URL",
			"user":           "user-id",
		},
	}
	v, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("error marshalling body: %s", err)
	}
	request := events.APIGatewayProxyRequest{
		Body: string(v),
	}
	body := ParseApiRequest(request)
	if body["action"] != float64(0) {
		t.Fatalf("action should be 0 but got %d", body["action"])
	}
	if body["params"].(map[string]interface{})["image"] != "https://aimless.ai/images/ai-canvas-logo.png" {
		t.Fatalf("image should be https://aimless.ai/images/ai-canvas-logo.png but got %s", body["params"].(map[string]interface{})["image"])
	}
	if body["params"].(map[string]interface{})["size"] != "512x512" {
		t.Fatalf("size should be 512x512 but got %s", body["params"].(map[string]interface{})["size"])
	}
	if body["params"].(map[string]interface{})["prompt"] != "something simple" {
		t.Fatalf("prompt should be something simple but got %s", body["params"].(map[string]interface{})["prompt"])
	}
	if body["params"].(map[string]interface{})["n"] != float64(1) {
		t.Fatalf("n should be 1 but got %d", body["params"].(map[string]interface{})["n"])
	}
	if body["params"].(map[string]interface{})["responseFormat"] != "URL" {
		t.Fatalf("responseFormat should be URL but got %s", body["params"].(map[string]interface{})["responseFormat"])
	}
	if body["params"].(map[string]interface{})["user"] != "user-id" {
		t.Fatalf("user should be user-id but got %s", body["params"].(map[string]interface{})["user"])
	}
}

func TestParseRequestAction(t *testing.T) {
	body := map[string]interface{}{
		"action": 0,
		"params": map[string]interface{}{
			"image":          "https://aimless.ai/images/ai-canvas-logo.png",
			"size":           "512x512",
			"prompt":         "something simple",
			"n":              1,
			"responseFormat": "URL",
			"user":           "user-id",
		},
	}
	action := ParseRequestAction(body)
	if action != aiTypes.GenerateImageAction {
		t.Fatalf("action should be GenerateImageAction but got %d", action)
	}
}

func TestConvertFloatToInt(t *testing.T) {
	i := ConvertFloatToInt(float64(1))
	if i != 1 {
		t.Fatalf("i should be 1 but got %d", i)
	}
	i = ConvertFloatToInt(1)
	if i != 1 {
		t.Fatalf("i should be 1 but got %d", i)
	}
	i = ConvertFloatToInt("1")
	if i != 1 {
		t.Fatalf("i should be 1 but got %d", i)
	}
}

func TestApiResponse(t *testing.T) {
	ai := aiTypes.AnalyticsItem{
		Id:      "1234",
		Success: true,
		Record: aiTypes.QueueRequest{
			Id: "1234",
		},
		Attempts: map[string]aiTypes.ImageResponseWrapper{
			"1": {
				Response: openai.ImageResponse{
					Created: 1234,
					Data: []openai.ImageResponseDataInner{
						{
							URL: "https://aimless.ai/images/ai-canvas-logo.png",
						},
					},
				},
			},
		},
	}
	response, err := ApiResponse(ai)
	if err != nil {
		t.Fatalf("error should be nil but got %s", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status code should be 200 but got %d", response.StatusCode)
	}
	if response.Body != "{\"id\":\"1234\",\"message\":\"Successfully added 1234\"}" {
		t.Fatalf("body should be {\"id\":\"1234\",\"message\":\"Successfully added 1234\"} but got %s", response.Body)
	}
}
