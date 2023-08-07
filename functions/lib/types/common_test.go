package types

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestGIRequestAction(t *testing.T) {
	var action RequestAction = GenerateImageAction
	if action != 0 {
		t.Fatalf("GenerateImageAction should be 1 but got %d", action)
	}
}

func TestEditImageAction(t *testing.T) {
	var action RequestAction = EditImageAction
	if action != 1 {
		t.Fatalf("EditImageAction should be 1 but got %d", action)
	}
}

func TestVariateImageAction(t *testing.T) {
	var action RequestAction = VariateImageAction
	if action != 2 {
		t.Fatalf("VariateImageAction should be 1 but got %d", action)
	}
}

// write a unit test for the struct ResponseFormat
func TestUrlResponseFormat(t *testing.T) {
	var urlFormat ResponseFormat
	urlFormat = "URL"
	if urlFormat.OpenaiResponseFormat() != "url" {
		t.Fatalf("jsonFormat should be json but got %s", urlFormat)
	}
}

func TestB64ResponseFormat(t *testing.T) {
	var b64Format ResponseFormat
	b64Format = "BASE64"
	if b64Format.OpenaiResponseFormat() != "b64_json" {
		t.Fatalf("jsonFormat should be json but got %s", b64Format)
	}
}

// write a unit test for the struct ImageSize
func TestSmallImageSize(t *testing.T) {
	var smallImageSize ImageSize
	smallImageSize = "256x256"
	if openai.CreateImageSize256x256 != smallImageSize {
		t.Fatalf("smallImageSize should be 256x256 but got %s", smallImageSize)
	}
	if smallImageSize.OpenaiImageSize() != "256x256" {
		t.Fatalf("smallImageSize should be 256x256 but got %s", smallImageSize)
	}
}

// write a unit test for medium ImageSize
func TestMediumImageSize(t *testing.T) {
	var mediumImageSize ImageSize
	mediumImageSize = "512x512"
	if openai.CreateImageSize512x512 != mediumImageSize {
		t.Fatalf("mediumImageSize should be 512x512 but got %s", mediumImageSize)
	}
	if mediumImageSize.OpenaiImageSize() != "512x512" {
		t.Fatalf("mediumImageSize should be 512x512 but got %s", mediumImageSize)
	}
}

// write a unit test for large ImageSize
func TestLargeImageSize(t *testing.T) {
	var largeImageSize ImageSize
	largeImageSize = "1024x1024"
	if openai.CreateImageSize1024x1024 != largeImageSize {
		t.Fatalf("largeImageSize should be 1024x1024 but got %s", largeImageSize)
	}
	if largeImageSize.OpenaiImageSize() != "1024x1024" {
		t.Fatalf("largeImageSize should be 1024x1024 but got %s", largeImageSize)
	}
}

// write a unit test for ImageSize with invalid value
func TestInvalidImageSize(t *testing.T) {
	var invalidImageSize ImageSize
	invalidImageSize = "invalid"
	if invalidImageSize.OpenaiImageSize() != "256x256" {
		t.Fatalf("invalidImageSize should be empty but got %s", invalidImageSize.OpenaiImageSize())
	}
}
