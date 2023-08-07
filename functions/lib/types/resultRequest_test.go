package types

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sashabaranov/go-openai"
)

func TestResultRequest(t *testing.T) {
	rr := ResultRequest{
		Record: QueueRequest{
			Metadata: CommonMetadata{
				TraceId: "1234",
			},
			Id:       "1234",
			Action:   EditImageAction,
			Priority: 1,
			CreateImageEdit: EditImageRequest{
				Prompt:         "testing",
				N:              1,
				Image:          "test",
				Mask:           "test",
				Size:           "512x512",
				User:           "1234",
				ResponseFormat: "BASE64",
			},
		},
		Result: ImageResponseWrapper{
			Success: true,
			Response: openai.ImageResponse{
				Created: 1234,
				Data: []openai.ImageResponseDataInner{
					{
						B64JSON: "test",
					},
				},
			},
		},
	}

	if rr.Record.Metadata.TraceId != "1234" {
		t.Fatalf("TraceId should be 1234 but got %s", rr.Record.Metadata.TraceId)
	}
	if rr.Record.Id != "1234" {
		t.Fatalf("Id should be 1234 but got %s", rr.Record.Id)
	}
	if rr.Record.Action != 1 {
		t.Fatalf("Action should be GenerateImageAction but got %d", rr.Record.Action)
	}
	if rr.Record.Priority != 1 {
		t.Fatalf("Priority should be 1 but got %d", rr.Record.Priority)
	}
	if rr.Record.CreateImageEdit.Image != "test" {
		t.Fatalf("Image should be test but got %s", rr.Record.CreateImageEdit.Image)
	}
	if rr.Record.CreateImageEdit.Mask != "test" {
		t.Fatalf("Mask should be test but got %s", rr.Record.CreateImageEdit.Mask)
	}
	if rr.Record.CreateImageEdit.N != 1 {
		t.Fatalf("N should be 1 but got %d", rr.Record.CreateImageEdit.N)
	}
	if rr.Record.CreateImageEdit.Size.OpenaiImageSize() != "512x512" {
		t.Fatalf("Size should be 512x512 but got %s", rr.Record.CreateImageEdit.Size)
	}
	if rr.Record.CreateImageEdit.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
		t.Fatalf("ResponseFormat should be b64_json but got %s", rr.Record.CreateImageEdit.ResponseFormat)
	}
	if rr.Record.CreateImageEdit.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", rr.Record.CreateImageEdit.User)
	}
	if rr.Result.Success != true {
		t.Fatalf("Success should be true but got %t", rr.Result.Success)
	}
	if rr.Result.Response.Created != 1234 {
		t.Fatalf("Created should be 1234 but got %d", rr.Result.Response.Created)
	}
	if rr.Result.Response.Data[0].B64JSON != "test" {
		t.Fatalf("B64JSON should be test but got %s", rr.Result.Response.Data[0].B64JSON)
	}

}

func TestResultRequestToDynamoDB(t *testing.T) {
	rr := ResultRequest{
		Record: QueueRequest{
			Metadata: CommonMetadata{
				TraceId: "1234",
			},
			Id:       "1234",
			Action:   EditImageAction,
			Priority: 1,
			CreateImageEdit: EditImageRequest{
				Prompt:         "testing",
				N:              1,
				Image:          "test",
				Mask:           "test",
				Size:           "512x512",
				User:           "1234",
				ResponseFormat: "BASE64",
			},
		},
		Result: ImageResponseWrapper{
			Success: true,
			Response: openai.ImageResponse{
				Created: 1234,
				Data: []openai.ImageResponseDataInner{
					{
						B64JSON: "test",
					},
				},
			},
		},
	}
	ddbItem := rr.ToDynamoDB()
	if ddbItem["id"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("id should be 1234 but got %s", ddbItem["id"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["priority"].(*types.AttributeValueMemberN).Value != "1" {
		t.Fatalf("priority should be 1 but got %s", ddbItem["priority"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["action"].(*types.AttributeValueMemberN).Value != "1" {
		t.Fatalf("action should be 1 but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["action"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["image"].(*types.AttributeValueMemberS).Value != "test" {
		t.Fatalf("image should be test but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["image"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["mask"].(*types.AttributeValueMemberS).Value != "test" {
		t.Fatalf("mask should be test but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["mask"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["n"].(*types.AttributeValueMemberN).Value != "1" {
		t.Fatalf("n should be 1 but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["n"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["size"].(*types.AttributeValueMemberS).Value != "512x512" {
		t.Fatalf("size should be 512x512 but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["size"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["user"].(*types.AttributeValueMemberS).Value != "1234" {
		t.Fatalf("user should be 1234 but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["user"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["responseFormat"].(*types.AttributeValueMemberS).Value != "BASE64" {
		t.Fatalf("responseFormat should be b64_json but got %s", ddbItem["request"].(*types.AttributeValueMemberM).Value["createImageEdit"].(*types.AttributeValueMemberM).Value["responseFormat"].(*types.AttributeValueMemberS).Value)
	}
	if ddbItem["result"].(*types.AttributeValueMemberM).Value["success"].(*types.AttributeValueMemberBOOL).Value != true {
		t.Fatalf("success should be true but got %t", ddbItem["result"].(*types.AttributeValueMemberM).Value["success"].(*types.AttributeValueMemberBOOL).Value)
	}
	if ddbItem["result"].(*types.AttributeValueMemberM).Value["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value != "1234" {
		t.Fatalf("created should be 1234 but got %s", ddbItem["result"].(*types.AttributeValueMemberM).Value["response"].(*types.AttributeValueMemberM).Value["created"].(*types.AttributeValueMemberN).Value)
	}
	if ddbItem["result"].(*types.AttributeValueMemberM).Value["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[0].(*types.AttributeValueMemberM).Value["b64"].(*types.AttributeValueMemberS).Value != "test" {
		t.Fatalf("b64_json should be test but got %s", ddbItem["result"].(*types.AttributeValueMemberM).Value["response"].(*types.AttributeValueMemberM).Value["data"].(*types.AttributeValueMemberL).Value[0].(*types.AttributeValueMemberM).Value["b64"].(*types.AttributeValueMemberS).Value)
	}

}

func TestResultRequestFromDynamoDB(t *testing.T) {
	for i := 0; i < 3; i++ {
		ddbItem := map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: "1234",
			},
			"priority": &types.AttributeValueMemberN{
				Value: "1",
			},
			"request": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"action": &types.AttributeValueMemberN{
						Value: fmt.Sprint(i),
					},
					"metadata": &types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"traceId": &types.AttributeValueMemberS{
								Value: "1234",
							},
						},
					},
					"createImageEdit": &types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"image": &types.AttributeValueMemberS{
								Value: "test",
							},
							"mask": &types.AttributeValueMemberS{
								Value: "test",
							},
							"prompt": &types.AttributeValueMemberS{
								Value: "testing",
							},
							"n": &types.AttributeValueMemberN{
								Value: "1",
							},
							"size": &types.AttributeValueMemberS{
								Value: "512x512",
							},
							"user": &types.AttributeValueMemberS{
								Value: "1234",
							},
							"responseFormat": &types.AttributeValueMemberS{
								Value: "BASE64",
							},
						},
					},
					"createImage": &types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"prompt": &types.AttributeValueMemberS{
								Value: "testing",
							},
							"n": &types.AttributeValueMemberN{
								Value: "1",
							},
							"size": &types.AttributeValueMemberS{
								Value: "512x512",
							},
							"user": &types.AttributeValueMemberS{
								Value: "1234",
							},
							"responseFormat": &types.AttributeValueMemberS{
								Value: "URL",
							},
						},
					},
					"createImageVariation": &types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"image": &types.AttributeValueMemberS{
								Value: "test",
							},
							"n": &types.AttributeValueMemberN{
								Value: "1",
							},
							"size": &types.AttributeValueMemberS{
								Value: "512x512",
							},
							"user": &types.AttributeValueMemberS{
								Value: "1234",
							},
							"responseFormat": &types.AttributeValueMemberS{
								Value: "BASE64",
							},
						},
					},
				},
			},
			"result": &types.AttributeValueMemberM{
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
											"b64": &types.AttributeValueMemberS{
												Value: "test",
											},
										},
									},
									&types.AttributeValueMemberM{
										Value: map[string]types.AttributeValue{
											"url": &types.AttributeValueMemberS{
												Value: "test",
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
		rr := ResultRequest{}
		rr.FromDynamoDB(ddbItem)
		if rr.Record.Metadata.TraceId != "1234" {
			t.Fatalf("TraceId should be 1234 but got %s", rr.Record.Metadata.TraceId)
		}
		if rr.Record.Id != "1234" {
			t.Fatalf("Id should be 1234 but got %s", rr.Record.Id)
		}
		if rr.Record.Action != RequestAction(i) {
			t.Fatalf("Action should be %d but got %d", i, rr.Record.Action)
		}
		if rr.Record.Priority != 1 {
			t.Fatalf("Priority should be 1 but got %d", rr.Record.Priority)
		}
		if i == 1 {

			if rr.Record.CreateImageEdit.Image != "test" {
				t.Fatalf("Image should be test but got %s", rr.Record.CreateImageEdit.Image)
			}
			if rr.Record.CreateImageEdit.Mask != "test" {
				t.Fatalf("Mask should be test but got %s", rr.Record.CreateImageEdit.Mask)
			}
			if rr.Record.CreateImageEdit.N != 1 {
				t.Fatalf("N should be %d but got %d", i, rr.Record.CreateImageEdit.N)
			}
			if rr.Record.CreateImageEdit.Size.OpenaiImageSize() != "512x512" {
				t.Fatalf("Size should be 512x512 but got %s", rr.Record.CreateImageEdit.Size)
			}
			if rr.Record.CreateImageEdit.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
				t.Fatalf("ResponseFormat should be b64_json but got %s", rr.Record.CreateImageEdit.ResponseFormat)
			}
			if rr.Record.CreateImageEdit.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", rr.Record.CreateImageEdit.User)
			}
		}
		if i == 0 {

			if rr.Record.CreateImage.Prompt != "testing" {
				t.Fatalf("Prompt should be testing but got %s", rr.Record.CreateImage.Prompt)
			}
			if rr.Record.CreateImage.N != 1 {
				t.Fatalf("N should be 1 but got %d", rr.Record.CreateImage.N)
			}
			if rr.Record.CreateImage.Size.OpenaiImageSize() != "512x512" {
				t.Fatalf("Size should be 512x512 but got %s", rr.Record.CreateImage.Size)
			}
			if rr.Record.CreateImage.ResponseFormat.OpenaiResponseFormat() != "url" {
				t.Fatalf("ResponseFormat should be url but got %s", rr.Record.CreateImage.ResponseFormat)
			}
			if rr.Record.CreateImage.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", rr.Record.CreateImage.User)
			}
		}
		if i == 2 {

			if rr.Record.CreateImageVariation.Image != "test" {
				t.Fatalf("Image should be test but got %s", rr.Record.CreateImageVariation.Image)
			}
			if rr.Record.CreateImageVariation.N != 1 {
				t.Fatalf("N should be %d but got %d", i, rr.Record.CreateImageVariation.N)
			}
			if rr.Record.CreateImageVariation.Size.OpenaiImageSize() != "512x512" {
				t.Fatalf("Size should be 512x512 but got %s", rr.Record.CreateImageVariation.Size)
			}
			if rr.Record.CreateImageVariation.ResponseFormat.OpenaiResponseFormat() != "b64_json" {
				t.Fatalf("ResponseFormat should be b64_json but got %s", rr.Record.CreateImageVariation.ResponseFormat)
			}
			if rr.Record.CreateImageVariation.User != "1234" {
				t.Fatalf("User should be 1234 but got %s", rr.Record.CreateImageVariation.User)
			}

		}
		if rr.Result.Success != true {
			t.Fatalf("Success should be true but got %t", rr.Result.Success)
		}
		if rr.Result.Response.Created != 1234 {
			t.Fatalf("Created should be 1234 but got %d", rr.Result.Response.Created)
		}
		if rr.Result.Response.Data[0].B64JSON != "test" {
			t.Fatalf("B64JSON should be test but got %s", rr.Result.Response.Data[0].B64JSON)
		}

	}

}
