package ai

import (
	"context"
	"io"
	"log"
	"os"

	. "github.com/aimless-it/ai-canvas/functions/lib/process"
	. "github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	openai "github.com/sashabaranov/go-openai"
)

// Environment variables required for this function:
// OPENAI_API_KEY_SECRET_ID - The AWS Secrets Manager secret ID that contains your OpenAI API key

func openaiConfig() openai.Client {
	secretName := os.Getenv("OPENAI_API_KEY_SECRET_ID")
	region := Region()

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString

	// Your OpenAI API key

	// Create a client
	client := openai.NewClient(secretString)

	return *client
}

func GenerateImage(request GenerateImageRequest) (openai.ImageResponse, error) {
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
	return client.CreateImage(context.Background(), imageRequest)
}

func EditImage(request EditImageRequest) (openai.ImageResponse, error) {
	// Create a client
	client := openaiConfig()
	editImageRequest := openai.ImageEditRequest{
		Prompt:         request.Prompt,
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.OpenaiResponseFormat(),
		Image:          GetImageFromS3(request.Image),
		Mask:           GetImageFromS3(request.Mask),
	}
	// Create the completion
	return client.CreateEditImage(context.Background(), editImageRequest)
}

func CreateImageVariation(request CreateImageVariantRequest) (openai.ImageResponse, error) {
	// Create a client
	client := openaiConfig()

	variantImageRequest := openai.ImageVariRequest{
		N:              request.N,
		Size:           request.Size.OpenaiImageSize(),
		ResponseFormat: request.ResponseFormat.OpenaiResponseFormat(),
		Image:          GetImageFromS3(request.Image),
	}
	// Create the completion
	return client.CreateVariImage(context.Background(), variantImageRequest)
}

func GetImageFromS3(imageName string) *os.File {
	region := Region()
	// create s3 client
	config, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
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
