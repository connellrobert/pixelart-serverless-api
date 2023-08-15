package internal

import (
	"testing"

	aiTypes "github.com/aimless-it/ai-canvas/functions/lib/types"
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
