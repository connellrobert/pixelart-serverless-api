package types

import "testing"

func TestGenerateImageStructure(t *testing.T) {
	genImage := GenerateImageRequest{
		Prompt:         "This is a test prompt",
		N:              1,
		Size:           "256x256",
		ResponseFormat: "URL",
		User:           "1234",
	}
	if genImage.Prompt != "This is a test prompt" {
		t.Fatalf("Prompt should be This is a test prompt but got %s", genImage.Prompt)
	}
	if genImage.N != 1 {
		t.Fatalf("N should be 1 but got %d", genImage.N)
	}
	if genImage.Size != "256x256" {
		t.Fatalf("Size should be 256x256 but got %s", genImage.Size)
	}
	if genImage.ResponseFormat != "URL" {
		t.Fatalf("ResponseFormat should be URL but got %s", genImage.ResponseFormat)
	}
	if genImage.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", genImage.User)
	}
}

func TestGenerateImageStructureWithDefaults(t *testing.T) {
	genImage := GenerateImageRequest{
		Prompt: "This is a test prompt",
	}
	if genImage.Prompt != "This is a test prompt" {
		t.Fatalf("Prompt should be This is a test prompt but got %s", genImage.Prompt)
	}
	//TODO: Should use default values when constructing structs
	if genImage.N != 0 {
		t.Fatalf("N should be 1 but got %d", genImage.N)
	}
	if genImage.Size.OpenaiImageSize() != "256x256" {
		t.Fatalf("Size should be 256x256 but got %s", genImage.Size)
	}
	if genImage.ResponseFormat.OpenaiResponseFormat() != "url" {
		t.Fatalf("ResponseFormat should be URL but got %s", genImage.ResponseFormat)
	}
	if genImage.User != "" {
		t.Fatalf("User should be empty but got %s", genImage.User)
	}
}

func TestEditImage(t *testing.T) {
	editImage := EditImageRequest{
		Prompt:         "This is a test prompt",
		N:              1,
		Size:           "512x512",
		ResponseFormat: "URL",
		User:           "1234",
		Mask:           "test",
		Image:          "test",
	}
	if editImage.Prompt != "This is a test prompt" {
		t.Fatalf("Prompt should be This is a test prompt but got %s", editImage.Prompt)
	}
	if editImage.N != 1 {
		t.Fatalf("N should be 1 but got %d", editImage.N)
	}
	if editImage.Size.OpenaiImageSize() != "512x512" {
		t.Fatalf("Size should be 256x256 but got %s", editImage.Size)
	}
	if editImage.ResponseFormat.OpenaiResponseFormat() != "url" {
		t.Fatalf("ResponseFormat should be URL but got %s", editImage.ResponseFormat)
	}
	if editImage.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", editImage.User)
	}
	if editImage.Mask != "test" {
		t.Fatalf("Mask should be test but got %s", editImage.Mask)
	}
	if editImage.Image != "test" {
		t.Fatalf("Image should be test but got %s", editImage.Image)
	}

}

// write a unit test for CreateImageVariantRequest struct
func TestCreateImageVariantRequest(t *testing.T) {
	variant := CreateImageVariantRequest{
		N:              1,
		Size:           "512x512",
		Image:          "test",
		ResponseFormat: "URL",
		User:           "1234",
	}
	if variant.N != 1 {
		t.Fatalf("N should be 1 but got %d", variant.N)
	}
	if variant.Size.OpenaiImageSize() != "512x512" {
		t.Fatalf("Size should be 256x256 but got %s", variant.Size)
	}
	if variant.ResponseFormat.OpenaiResponseFormat() != "url" {
		t.Fatalf("ResponseFormat should be URL but got %s", variant.ResponseFormat)
	}
	if variant.User != "1234" {
		t.Fatalf("User should be 1234 but got %s", variant.User)
	}
	if variant.Image != "test" {
		t.Fatalf("Image should be test but got %s", variant.Image)
	}

}
