package types

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestQueueRequestFromDynamoDB(t *testing.T) {
	prompt := "This is a test prompt"
	for i := 0; i < 3; i++ {

		ddbItem := map[string]types.AttributeValue{
			"action": &types.AttributeValueMemberN{
				Value: fmt.Sprint(i),
			},
			"createImage": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"prompt": &types.AttributeValueMemberS{
						Value: prompt,
					},
					"n": &types.AttributeValueMemberN{
						Value: "1",
					},
					"size": &types.AttributeValueMemberS{
						Value: "256x256",
					},
					"responseFormat": &types.AttributeValueMemberS{
						Value: "URL",
					},
					"user": &types.AttributeValueMemberS{
						Value: "1234",
					},
				},
			},
			"createImageEdit": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"prompt": &types.AttributeValueMemberS{
						Value: prompt,
					},
					"n": &types.AttributeValueMemberN{
						Value: "1",
					},
					"size": &types.AttributeValueMemberS{
						Value: "1024x1024",
					},
					"responseFormat": &types.AttributeValueMemberS{
						Value: "URL",
					},
					"user": &types.AttributeValueMemberS{
						Value: "1234",
					},
					"image": &types.AttributeValueMemberS{
						Value: "test",
					},
					"mask": &types.AttributeValueMemberS{
						Value: "test",
					},
				},
			},
			"createImageVariation": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"n": &types.AttributeValueMemberN{
						Value: "1",
					},
					"size": &types.AttributeValueMemberS{
						Value: "512x512",
					},
					"responseFormat": &types.AttributeValueMemberS{
						Value: "URL",
					},
					"user": &types.AttributeValueMemberS{
						Value: "1234",
					},
					"image": &types.AttributeValueMemberS{
						Value: "test",
					},
				},
			},
			"metadata": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"traceId": &types.AttributeValueMemberS{
						Value: "1234",
					},
				},
			},
			"id": &types.AttributeValueMemberS{
				Value: "1234",
			},
			"priority": &types.AttributeValueMemberN{
				Value: "1",
			},
		}
		var queueRequest QueueRequest
		queueRequest.FromDynamoDB(ddbItem)
		if queueRequest.Action != RequestAction(i) {
			t.Fatalf("Action should be %d but got %d", i, queueRequest.Action)
		}
		if queueRequest.Metadata.TraceId != "1234" {
			t.Fatalf("TraceId should be 1234 but got %s", queueRequest.Metadata.TraceId)
		}
		if queueRequest.Id != "1234" {
			t.Fatalf("Id should be 1234 but got %s", queueRequest.Id)
		}
		if queueRequest.Priority != 1 {
			t.Fatalf("Priority should be 1 but got %d", queueRequest.Priority)
		}
		switch queueRequest.Action {
		case GenerateImageAction:
			if queueRequest.CreateImage.Prompt != prompt {
				t.Fatalf("Prompt should be %s but got %s", prompt, queueRequest.CreateImage.Prompt)
			}
			if queueRequest.CreateImage.N != 1 {
				t.Fatalf("N should be 1 but got %d", queueRequest.CreateImage.N)
			}
			if queueRequest.CreateImage.Size != "256x256" {
				t.Fatalf("Size should be 256x256 but got %s", queueRequest.CreateImage.Size)
			}
			if queueRequest.CreateImage.ResponseFormat != "URL" {
				t.Fatalf("ResponseFormat should be URL but got %s", queueRequest.CreateImage.ResponseFormat)
			}
			if queueRequest.CreateImage.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", queueRequest.CreateImage.User)
			}
		case EditImageAction:
			if queueRequest.CreateImageEdit.Prompt != prompt {
				t.Fatalf("Prompt should be %s but got %s", prompt, queueRequest.CreateImageEdit.Prompt)
			}
			if queueRequest.CreateImageEdit.N != 1 {
				t.Fatalf("N should be 1 but got %d", queueRequest.CreateImageEdit.N)
			}
			if queueRequest.CreateImageEdit.Size != "1024x1024" {
				t.Fatalf("Size should be 1024x1024 but got %s", queueRequest.CreateImageEdit.Size)
			}
			if queueRequest.CreateImageEdit.ResponseFormat != "URL" {
				t.Fatalf("ResponseFormat should be URL but got %s", queueRequest.CreateImageEdit.ResponseFormat)
			}
			if queueRequest.CreateImageEdit.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", queueRequest.CreateImageEdit.User)
			}
			if queueRequest.CreateImageEdit.Image != "test" {
				t.Fatalf("Image should be test but got %s", queueRequest.CreateImageEdit.Image)
			}
			if queueRequest.CreateImageEdit.Mask != "test" {
				t.Fatalf("Mask should be test but got %s", queueRequest.CreateImageEdit.Mask)
			}
		case VariateImageAction:
			if queueRequest.CreateImageVariation.N != 1 {
				t.Fatalf("N should be 1 but got %d", queueRequest.CreateImageVariation.N)
			}
			if queueRequest.CreateImageVariation.Size != "512x512" {
				t.Fatalf("Size should be 512x512 but got %s", queueRequest.CreateImageVariation.Size)
			}
			if queueRequest.CreateImageVariation.ResponseFormat != "URL" {
				t.Fatalf("ResponseFormat should be URL but got %s", queueRequest.CreateImageVariation.ResponseFormat)
			}
			if queueRequest.CreateImageVariation.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", queueRequest.CreateImageVariation.User)
			}
			if queueRequest.CreateImageVariation.Image != "test" {
				t.Fatalf("Image should be test but got %s", queueRequest.CreateImageVariation.Image)
			}

		}
	}
}

func TestMapParams(t *testing.T) {
	cases := []RequestAction{GenerateImageAction, EditImageAction, VariateImageAction}
	params := map[string]interface{}{
		"prompt":         "This is a test prompt",
		"n":              1.0,
		"size":           "256x256",
		"responseFormat": "URL",
		"user":           "1234",
		"image":          "test",
		"mask":           "test",
	}
	for _, c := range cases {
		qr := QueueRequest{}
		qr.MapParams(c, params)
		switch c {
		case GenerateImageAction:
			if qr.CreateImage.Prompt != "This is a test prompt" {
				t.Fatalf("Prompt should be This is a test prompt but got %s", qr.CreateImage.Prompt)
			}
			if qr.CreateImage.N != 1 {
				t.Fatalf("N should be 1 but got %d", qr.CreateImage.N)
			}
			if qr.CreateImage.Size != "256x256" {
				t.Fatalf("Size should be 256x256 but got %s", qr.CreateImage.Size)
			}
			if qr.CreateImage.ResponseFormat != "URL" {
				t.Fatalf("ResponseFormat should be URL but got %s", qr.CreateImage.ResponseFormat)
			}
			if qr.CreateImage.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", qr.CreateImage.User)
			}
		case EditImageAction:
			if qr.CreateImageEdit.Prompt != "This is a test prompt" {
				t.Fatalf("Prompt should be This is a test prompt but got %s", qr.CreateImageEdit.Prompt)
			}
			if qr.CreateImageEdit.N != 1 {
				t.Fatalf("N should be 1 but got %d", qr.CreateImageEdit.N)
			}
			if qr.CreateImageEdit.Size != "256x256" {
				t.Fatalf("Size should be 256x256 but got %s", qr.CreateImageEdit.Size)
			}
			if qr.CreateImageEdit.ResponseFormat != "URL" {
				t.Fatalf("ResponseFormat should be URL but got %s", qr.CreateImageEdit.ResponseFormat)
			}
			if qr.CreateImageEdit.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", qr.CreateImageEdit.User)
			}
			if qr.CreateImageEdit.Image != "test" {
				t.Fatalf("Image should be test but got %s", qr.CreateImageEdit.Image)
			}
			if qr.CreateImageEdit.Mask != "test" {
				t.Fatalf("Mask should be test but got %s", qr.CreateImageEdit.Mask)
			}
		case VariateImageAction:
			if qr.CreateImageVariation.N != 1 {
				t.Fatalf("N should be 1 but got %d", qr.CreateImageVariation.N)
			}
			if qr.CreateImageVariation.Size != "256x256" {
				t.Fatalf("Size should be 256x256 but got %s", qr.CreateImageVariation.Size)
			}
			if qr.CreateImageVariation.ResponseFormat != "URL" {
				t.Fatalf("ResponseFormat should be URL but got %s", qr.CreateImageVariation.ResponseFormat)
			}
			if qr.CreateImageVariation.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", qr.CreateImageVariation.User)
			}
			if qr.CreateImageVariation.Image != "test" {
				t.Fatalf("Image should be test but got %s", qr.CreateImageVariation.Image)
			}
		}
	}
}
