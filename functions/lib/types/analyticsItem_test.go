package types

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	openai "github.com/sashabaranov/go-openai"
)

func TestAnalyticsItem(t *testing.T) {
	ai := AnalyticsItem{
		Success: true,
		Id:      "1234",
		Record: QueueRequest{
			Metadata: CommonMetadata{
				TraceId: "1234",
			},
			Id:       "1234",
			Action:   GenerateImageAction,
			Priority: 1,
			CreateImage: GenerateImageRequest{
				Prompt:         "This is a test prompt",
				N:              1,
				Size:           "1024x1024",
				ResponseFormat: "BASE64",
				User:           "1234",
			},
			CreateImageEdit:      EditImageRequest{},
			CreateImageVariation: CreateImageVariantRequest{},
		},
		Attempts: map[string]ImageResponseWrapper{
			"1": {
				Success: true,
				Response: openai.ImageResponse{
					Created: 1234,
					Data: []openai.ImageResponseDataInner{
						{
							URL: "https://test.com",
						},
					},
				},
			},
		},
	}
	if ai.Success != true {
		t.Fatalf("Success should be true but got %t", ai.Success)
	}
	if ai.Id != "1234" {
		t.Fatalf("Id should be 1234 but got %s", ai.Id)
	}
	if ai.Record.Metadata.TraceId != "1234" {
		t.Fatalf("TraceId should be 1234 but got %s", ai.Record.Metadata.TraceId)
	}
	if ai.Record.Id != "1234" {
		t.Fatalf("Id should be 1234 but got %s", ai.Record.Id)
	}
	if ai.Record.Action != 0 {
		t.Fatalf("Action should be GenerateImageAction but got %d", ai.Record.Action)
	}
	if ai.Record.Priority != 1 {
		t.Fatalf("Priority should be 1 but got %d", ai.Record.Priority)
	}
	if ai.Record.CreateImage.Prompt != "This is a test prompt" {
		t.Fatalf("Prompt should be This is a test prompt but got %s", ai.Record.CreateImage.Prompt)
	}
	if ai.Record.CreateImage.N != 1 {
		t.Fatalf("N should be 1 but got %d", ai.Record.CreateImage.N)
	}
	if ai.Record.CreateImage.Size.OpenaiImageSize() != "1024x1024" {
		t.Fatalf("Size should be 1024x1024 but got %s", ai.Record.CreateImage.Size)
	}
	if ai.Record.CreateImage.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
		t.Fatalf("ResponseFormat should be b64_json but got %s", ai.Record.CreateImage.ResponseFormat)
	}
	if ai.Record.CreateImage.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", ai.Record.CreateImage.User)
	}
	if ai.Attempts["1"].Success != true {
		t.Fatalf("Success should be true but got %t", ai.Attempts["1"].Success)
	}
	if ai.Attempts["1"].Response.Created != 1234 {
		t.Fatalf("Created should be 1234 but got %d", ai.Attempts["1"].Response.Created)
	}
	if ai.Attempts["1"].Response.Data[0].URL != "https://test.com" {
		t.Fatalf("URL should be https://test.com but got %s", ai.Attempts["1"].Response.Data[0].URL)
	}

}

// write a unit test for the `ToDynamoDB` function for the AnalyticsItem struct
func TestAnalyticsItemToDynamoDB(t *testing.T) {
	analyticsItem := AnalyticsItem{
		Success: true,
		Id:      "1234",
		Record: QueueRequest{
			Metadata: CommonMetadata{
				TraceId: "1234",
			},
			Id:       "1234",
			Action:   GenerateImageAction,
			Priority: 1,
			CreateImage: GenerateImageRequest{
				Prompt:         "This is a test prompt",
				N:              1,
				Size:           "1024x1024",
				ResponseFormat: "BASE64",
				User:           "1234",
			},
		},
		Attempts: map[string]ImageResponseWrapper{},
	}
	ddbItem := analyticsItem.ToDynamoDB()
	if ddbItem["success"].(*types.AttributeValueMemberBOOL).Value != true {
		t.Fatalf("Success should be true but got %t", ddbItem["success"].(*types.AttributeValueMemberBOOL).Value)
	}
	if ddbItem["id"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("Id should be 1234 but got %s", ddbItem["id"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["metadata"].(*types.AttributeValueMemberM).Value["traceId"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("TraceId should be 1234 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["metadata"].(*types.AttributeValueMemberM).Value["traceId"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["id"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("Id should be 1234 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["id"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["action"].(*types.AttributeValueMemberN).Value != "0" {
		t.Fatalf("Action should be 0 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["action"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["priority"].(*types.AttributeValueMemberN).Value != "1" {
		t.Fatalf("Priority should be 1 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["priority"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["prompt"].(*types.AttributeValueMemberS).Value != "This is a test prompt" {
		t.Fatalf("Prompt should be This is a test prompt but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["prompt"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["n"].(*types.AttributeValueMemberN).Value != "1" {
		t.Fatalf("N should be 1 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["n"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["size"].(*types.AttributeValueMemberS).Value != "1024x1024" {
		t.Fatalf("Size should be 1024x1024 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["size"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["responseFormat"].(*types.AttributeValueMemberS).Value != "BASE64" {
		t.Fatalf("ResponseFormat should be BASE64 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["responseFormat"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["user"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("User should be 1234 but got %s", ddbItem["record"].(*types.AttributeValueMemberM).Value["createImage"].(*types.AttributeValueMemberM).Value["user"].(*types.AttributeValueMemberS).Value)
	}

}

// write a unit test for the `FromDynamoDB` function for the AnalyticsItem struct
func TestAnalyticsItemFromDynamoDB(t *testing.T) {
	prompt := "This is a test prompt"
	for i := 0; i < 3; i++ {
		dbItem := map[string]types.AttributeValue{
			"success": &types.AttributeValueMemberBOOL{
				Value: true,
			},
			"id": &types.AttributeValueMemberS{
				Value: "1234",
			},
			"record": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
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
					"action": &types.AttributeValueMemberN{
						Value: fmt.Sprint(i),
					},
					"priority": &types.AttributeValueMemberN{
						Value: "1",
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
								Value: "BASE64",
							},
							"user": &types.AttributeValueMemberS{
								Value: "1234",
							},
							"mask": &types.AttributeValueMemberS{
								Value: "1234",
							},
							"image": &types.AttributeValueMemberS{
								Value: "1234",
							},
						},
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
								Value: "1024x1024",
							},
							"responseFormat": &types.AttributeValueMemberS{
								Value: "BASE64",
							},
							"user": &types.AttributeValueMemberS{
								Value: "1234",
							},
						},
					},
					"createImageVariation": &types.AttributeValueMemberM{
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
								Value: "BASE64",
							},
							"user": &types.AttributeValueMemberS{
								Value: "1234",
							},
							"image": &types.AttributeValueMemberS{
								Value: "1234",
							},
						},
					},
				},
			},
			"attempts": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{},
			},
		}
		analyticsItem := AnalyticsItem{}
		analyticsItem.FromDynamoDB(dbItem)
		if analyticsItem.Success != true {
			t.Fatalf("Success should be true but got %t", analyticsItem.Success)
		}
		if analyticsItem.Id != "1234" {
			t.Fatalf("Id should be 1234 but got %s", analyticsItem.Id)
		}
		if analyticsItem.Record.Metadata.TraceId != "1234" {
			t.Fatalf("TraceId should be 1234 but got %s", analyticsItem.Record.Metadata.TraceId)
		}
		if analyticsItem.Record.Id != "1234" {
			t.Fatalf("Id should be 1234 but got %s", analyticsItem.Record.Id)
		}
		if analyticsItem.Record.Action != RequestAction(i) {
			t.Fatalf("Action should be 1 but got %d", analyticsItem.Record.Action)
		}
		if analyticsItem.Record.Priority != 1 {
			t.Fatalf("Priority should be 1 but got %d", analyticsItem.Record.Priority)
		}
		if i == 1 {

			if analyticsItem.Record.CreateImageEdit.Prompt != prompt {
				t.Fatalf("Prompt should be This is a test prompt but got %s", analyticsItem.Record.CreateImage.Prompt)
			}
			if analyticsItem.Record.CreateImageEdit.N != 1 {
				t.Fatalf("N should be 1 but got %d", analyticsItem.Record.CreateImage.N)
			}
			if analyticsItem.Record.CreateImageEdit.Size.OpenaiImageSize() != "1024x1024" {
				t.Fatalf("Size should be 1024x1024 but got %s", analyticsItem.Record.CreateImage.Size)
			}
			if analyticsItem.Record.CreateImageEdit.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
				t.Fatalf("ResponseFormat should be BASE64 but got %s", analyticsItem.Record.CreateImage.ResponseFormat)
			}
			if analyticsItem.Record.CreateImageEdit.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", analyticsItem.Record.CreateImage.User)
			}
			if analyticsItem.Record.CreateImageEdit.Mask != "1234" {
				t.Fatalf("Mask should be 1234 but got %s", analyticsItem.Record.CreateImageEdit.Mask)
			}
			if analyticsItem.Record.CreateImageEdit.Image != "1234" {
				t.Fatalf("Image should be 1234 but got %s", analyticsItem.Record.CreateImageEdit.Image)
			}
		}
		if i == 0 {
			if analyticsItem.Record.CreateImage.Prompt != prompt {
				t.Fatalf("Prompt should be This is a test prompt but got %s", analyticsItem.Record.CreateImage.Prompt)
			}
			if analyticsItem.Record.CreateImage.N != 1 {
				t.Fatalf("N should be 1 but got %d", analyticsItem.Record.CreateImage.N)
			}
			if analyticsItem.Record.CreateImage.Size.OpenaiImageSize() != "1024x1024" {
				t.Fatalf("Size should be 1024x1024 but got %s", analyticsItem.Record.CreateImage.Size)
			}
			if analyticsItem.Record.CreateImage.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
				t.Fatalf("ResponseFormat should be BASE64 but got %s", analyticsItem.Record.CreateImage.ResponseFormat)
			}
			if analyticsItem.Record.CreateImage.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", analyticsItem.Record.CreateImage.User)
			}

		}
		if i == 2 {
			if analyticsItem.Record.CreateImageVariation.N != 1 {
				t.Fatalf("N should be 1 but got %d", analyticsItem.Record.CreateImage.N)
			}
			if analyticsItem.Record.CreateImageVariation.Size.OpenaiImageSize() != "1024x1024" {
				t.Fatalf("Size should be 1024x1024 but got %s", analyticsItem.Record.CreateImage.Size)
			}
			if analyticsItem.Record.CreateImageVariation.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
				t.Fatalf("ResponseFormat should be BASE64 but got %s", analyticsItem.Record.CreateImage.ResponseFormat)
			}
			if analyticsItem.Record.CreateImageVariation.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", analyticsItem.Record.CreateImage.User)
			}
		}

		if len(analyticsItem.Attempts) != 0 {
			t.Fatalf("Attempts should be empty but got %d", len(analyticsItem.Attempts))
		}
	}
}

func TestAnalyticsAttemptsToDB(t *testing.T) {
	ai := AnalyticsItem{
		Success: true,
		Id:      "1234",
		Record: QueueRequest{
			Metadata: CommonMetadata{
				TraceId: "1234",
			},
			Id:       "1234",
			Action:   GenerateImageAction,
			Priority: 1,
			CreateImage: GenerateImageRequest{
				Prompt:         "This is a test prompt",
				N:              1,
				Size:           "1024x1024",
				ResponseFormat: "BASE64",
				User:           "1234",
			},
			CreateImageEdit:      EditImageRequest{},
			CreateImageVariation: CreateImageVariantRequest{},
		},
		Attempts: map[string]ImageResponseWrapper{
			"1": {
				Success: true,
				Response: openai.ImageResponse{
					Created: 1234,
					Data: []openai.ImageResponseDataInner{
						{
							URL: "test.com",
						},
						{
							URL: "test2.com",
						},
					},
				},
			},
			"2": {
				Success: true,
				Response: openai.ImageResponse{
					Created: 1234,
					Data: []openai.ImageResponseDataInner{
						{
							URL: "test.com",
						},
						{
							URL: "test2.com",
						},
					},
				},
			},
		},
	}

	ddbItem := ai.AttemptsToDynamoDB()
	attempt1 := ddbItem["1"].(*types.AttributeValueMemberM).Value
	attempt2 := ddbItem["2"].(*types.AttributeValueMemberM).Value
	if attempt1["success"].(*types.AttributeValueMemberBOOL).Value != true {
		t.Fatalf("Success should be true but got %t", attempt1["success"].(*types.AttributeValueMemberBOOL).Value)
	}
	if attempt1["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value != "1234" {
		t.Fatalf("Created should be 1234 but got %s", attempt1["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value)
	}
	if attempt1["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[0].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value != "test.com" {
		t.Fatalf("URL should be test.com but got %s", attempt1["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[0].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value)
	}
	if attempt1["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[1].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value != "test2.com" {
		t.Fatalf("URL should be test2.com but got %s", attempt1["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[1].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value)
	}
	if attempt2["success"].(*types.AttributeValueMemberBOOL).Value != true {
		t.Fatalf("Success should be true but got %t", attempt2["success"].(*types.AttributeValueMemberBOOL).Value)
	}
	if attempt2["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value != "1234" {
		t.Fatalf("Created should be 1234 but got %s", attempt2["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value)
	}
	if attempt2["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[0].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value != "test.com" {
		t.Fatalf("URL should be test.com but got %s", attempt2["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[0].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value)
	}
	if attempt2["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[1].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value != "test2.com" {
		t.Fatalf("URL should be test2.com but got %s", attempt2["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[1].(*types.AttributeValueMemberM).Value["url"].(*types.AttributeValueMemberS).Value)
	}

}

func TestAnalyticsAttemptsFromDB(t *testing.T) {
	aidb := map[string]types.AttributeValue{
		"1": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"success": &types.AttributeValueMemberBOOL{
					Value: true,
				},
				"response": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"created": &types.AttributeValueMemberN{
							Value: "1234",
						},
						"data": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"url": &types.AttributeValueMemberS{
											Value: "test.com",
										},
									},
								},
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"b64": &types.AttributeValueMemberS{
											Value: "dGVzdDIuY29tCg==",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"2": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"success": &types.AttributeValueMemberBOOL{
					Value: true,
				},
				"response": &types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"created": &types.AttributeValueMemberN{
							Value: "1234",
						},
						"data": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"url": &types.AttributeValueMemberS{
											Value: "test.com",
										},
									},
								},
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"b64": &types.AttributeValueMemberS{
											Value: "dGVzdDIuY29tCg==",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	analyticsItem := AnalyticsItem{}
	analyticsItem.AttemptsFromDynamoDB(aidb)
	if len(analyticsItem.Attempts) != 2 {
		t.Fatalf("Attempts should have 2 items but got %d", len(analyticsItem.Attempts))
	}
	attempt1 := analyticsItem.Attempts["1"]
	attempt2 := analyticsItem.Attempts["2"]
	if attempt1.Success != true {
		t.Fatalf("Success should be true but got %t", attempt1.Success)
	}
	if attempt1.Response.Created != 1234 {
		t.Fatalf("Created should be 1234 but got %d", attempt1.Response.Created)
	}
	if attempt1.Response.Data[0].URL != "test.com" {
		t.Fatalf("URL should be test.com but got %s", attempt1.Response.Data[0].URL)
	}
	if attempt1.Response.Data[1].B64JSON != "dGVzdDIuY29tCg==" {
		t.Fatalf("URL should be test2.com but got %s", attempt1.Response.Data[1].URL)
	}
	if attempt2.Success != true {
		t.Fatalf("Success should be true but got %t", attempt2.Success)
	}
	if attempt2.Response.Created != 1234 {
		t.Fatalf("Created should be 1234 but got %d", attempt2.Response.Created)
	}
	if attempt2.Response.Data[0].URL != "test.com" {
		t.Fatalf("URL should be test.com but got %s", attempt2.Response.Data[0].URL)
	}
	if attempt2.Response.Data[1].B64JSON != "dGVzdDIuY29tCg==" {
		t.Fatalf("URL should be test2.com but got %s", attempt2.Response.Data[1].URL)
	}

}
