package ai

import (
	"context"
	"io"
	"os"

	"github.com/connellrobert/pixelart-serverless-api/functions/lib/process"
	"github.com/connellrobert/pixelart-serverless-api/functions/lib/types"
	openai "github.com/sashabaranov/go-openai"
)

// Environment variables required for this function:
// OPENAI_API_KEY_SECRET_ID - The AWS Secrets Manager secret ID that contains your OpenAI API key

type aiProcessor interface {
	GenerateImage(client openai.Client, request openai.ImageRequest) (openai.ImageResponse, error)
	EditImage(client openai.Client, request openai.ImageEditRequest) (openai.ImageResponse, error)
	CreateImageVariation(client openai.Client, request openai.ImageVariRequest) (openai.ImageResponse, error)
	GetSecretValue(secretId string) string
	GetImageFromS3(image string) io.ReadCloser
	SaveFile(fileName string, fileContents io.ReadCloser) *os.File
}

type aiProcess struct{}

func (a aiProcess) GenerateImage(client openai.Client, request openai.ImageRequest) (openai.ImageResponse, error) {
	return client.CreateImage(context.Background(), request)
}

func (a aiProcess) EditImage(client openai.Client, request openai.ImageEditRequest) (openai.ImageResponse, error) {
	return client.CreateEditImage(context.Background(), request)
}

func (a aiProcess) CreateImageVariation(client openai.Client, request openai.ImageVariRequest) (openai.ImageResponse, error) {
	return client.CreateVariImage(context.Background(), request)
}

func (a aiProcess) GetSecretValue(secretId string) string {
	return process.GetSecretValue(secretId)
}

func (a aiProcess) GetImageFromS3(image string) io.ReadCloser {
	return process.GetImageFromS3(image)
}

func (a aiProcess) SaveFile(fileName string, fileContents io.ReadCloser) *os.File {
	return process.SaveFile(fileName, fileContents)
}

var ai aiProcessor

func init() {
	ai = aiProcess{}
}

func openaiConfig() openai.Client {
	secretString := ai.GetSecretValue("OPENAI_API_KEY_SECRET_ID")
	// Your OpenAI API key

	// Create a client
	client := openai.NewClient(secretString)

	return *client
}

func GenerateImage(request types.GenerateImageRequest) (openai.ImageResponse, error) {
	// Create a client
	client := openaiConfig()

	imageRequest := openai.ImageRequest{
		Prompt:         request.Prompt,
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.OpenaiResponseFormat(),
		User:           request.User,
	}

	// Create the completion
	return ai.GenerateImage(client, imageRequest)
}

func EditImage(request types.EditImageRequest) (openai.ImageResponse, error) {
	image := ai.SaveFile(request.Image, ai.GetImageFromS3(request.Image))
	mask := ai.SaveFile(request.Mask, ai.GetImageFromS3(request.Mask))
	// Create a client
	client := openaiConfig()
	editImageRequest := openai.ImageEditRequest{
		Prompt:         request.Prompt,
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.OpenaiResponseFormat(),
		Image:          image,
		Mask:           mask,
	}
	// Create the completion
	return ai.EditImage(client, editImageRequest)
}

func CreateImageVariation(request types.CreateImageVariantRequest) (openai.ImageResponse, error) {
	image := ai.SaveFile(request.Image, ai.GetImageFromS3(request.Image))
	// Create a client
	client := openaiConfig()

	variantImageRequest := openai.ImageVariRequest{
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.OpenaiResponseFormat(),
		Image:          image,
	}
	// Create the completion
	return ai.CreateImageVariation(client, variantImageRequest)
}
