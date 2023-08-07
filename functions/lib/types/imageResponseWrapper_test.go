package types

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestImageResponseWrapper(t *testing.T) {
	wrapper := ImageResponseWrapper{
		Success: true,
		Response: openai.ImageResponse{
			Created: 1234,
			Data: []openai.ImageResponseDataInner{
				{
					URL: "https://test.com",
				},
			},
		},
	}
	if wrapper.Success != true {
		t.Fatalf("Success should be true but got %t", wrapper.Success)
	}
	if wrapper.Response.Created != 1234 {
		t.Fatalf("Created should be 1234 but got %d", wrapper.Response.Created)
	}
	if wrapper.Response.Data[0].URL != "https://test.com" {
		t.Fatalf("URL should be https://test.com but got %s", wrapper.Response.Data[0].URL)
	}

}
