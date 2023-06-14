package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func SendResult(record interface{}, response interface{}) string {
	// Create a Lambda client
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		panic(err)
	}
	// invoke lambda
	lambdaArn := os.Getenv("RESULT_FUNCTION_ARN")
	svc := lambda.NewFromConfig(cfg)
	input := &lambda.InvokeInput{
		FunctionName: aws.String(lambdaArn),
		// payload should be a json string of the input arguments
		Payload: []byte(fmt.Sprintf(`{"record": %v, "response": %v}`, record, response)),
	}
	result, err := svc.Invoke(context.Background(), input)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result.Payload))
	return string(result.Payload)
}
