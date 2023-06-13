package main

import (
	"context"
	"os"

	"github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/config"
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
	apiKey := result.SecretString

	// Create a client
	client := openai.NewClient(apiKey)

	return client
}

func GenerateImage(request types.GenerateImageRequest) {
	// Create a client
	client := openaiConfig()

	// Create the completion
	completion, err := client.CreateImage(context.Background(), request)
	if err != nil {
		panic(err)
	}

	return completion.Data[0]
}
