package lib

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	openai "github.com/sashabaranov/go-openai"
)

// Environment variables required for this function:
// OPENAI_API_KEY_SECRET_ID - The AWS Secrets Manager secret ID that contains your OpenAI API key

func openaiConfig() openai.Client {
	// Create a Secrets Manager client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"),
	)
	if err != nil {
		panic(err)
	}
	svc := secretsmanager.NewFromConfig(cfg)

	openaiEnvVar := os.Getenv("OPENAI_API_KEY_SECRET_ID")
	// Get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(openaiEnvVar),
	}
	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	// Your OpenAI API key
	apiKey := *result.SecretString

	// Create a client
	client := openai.NewClient(apiKey)

	return *client
}

func GenerateImage(request GenerateImageRequest) openai.ImageResponse {
	// Create a client
	client := openaiConfig()

	imageRequest := openai.ImageRequest{
		Prompt:         request.Prompt,
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.openaiResponseFormat(),
		User:           request.User,
	}

	// Create the completion
	completion, err := client.CreateImage(context.Background(), imageRequest)
	if err != nil {
		panic(err)
	}

	return completion
}

func EditImage(request EditImageRequest) openai.ImageResponse {
	// Create a client
	client := openaiConfig()
	editImageRequest := openai.ImageEditRequest{
		Prompt:         request.Prompt,
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.openaiResponseFormat(),
		Image:          GetImageFromS3(request.Image),
		Mask:           GetImageFromS3(request.Mask),
	}
	// Create the completion
	completion, err := client.CreateEditImage(context.Background(), editImageRequest)
	if err != nil {
		panic(err)
	}

	return completion
}

func CreateImageVariation(request CreateImageVariantRequest) openai.ImageResponse {
	// Create a client
	client := openaiConfig()

	variantImageRequest := openai.ImageVariRequest{
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.openaiResponseFormat(),
		Image:          GetImageFromS3(request.Image),
	}
	// Create the completion
	completion, err := client.CreateVariImage(context.Background(), variantImageRequest)
	if err != nil {
		panic(err)
	}

	return completion
}

func GetImageFromS3(imageName string) *os.File {

	// create s3 client
	config, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	svc := s3.NewFromConfig(config)
	obj := &s3.GetObjectInput{
		Bucket: aws.String("openai-image-storage"),
		Key:    aws.String(imageName),
	}
	result, err := svc.GetObject(context.Background(), obj)
	if err != nil {
		panic(err)
	}
	file := os.NewFile(4, imageName)
	_, err = io.Copy(file, result.Body)
	if err != nil {
		panic(err)
	}
	return file
}
