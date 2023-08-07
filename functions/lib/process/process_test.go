package process

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/aimless-it/ai-canvas/functions/lib/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/stretchr/testify/mock"
)

func TestGetAWSConfig(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
	}()
	cfg := GetAWSConfig()
	if cfg.Region != region {
		t.Fatalf("Region should be %s but got %s", region, cfg.Region)
	}
}

func TestGetAwsConfigError(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(func(o *config.LoadOptions) error {
		return nil
	})
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
	}()
	GetAWSConfig()
}

func TestRegion(t *testing.T) {
	_region := Region()
	if _region != region {
		t.Fatalf("Region should be %s but got %s", region, _region)
	}
}

func TestRegionDefault(t *testing.T) {
	os.Unsetenv("AWS_REGION")
	_region := Region()
	if _region != "us-east-1" {
		t.Fatalf("Region should be us-east-1 but got %s", _region)
	}
}

func TestSendResult(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	sm := p.(*mockProcess).On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&sqs.SendMessageOutput{}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
		sm.Unset()
	}()
	SendResult(types.QueueRequest{}, types.ImageResponseWrapper{})
	if p.(*mockProcess).AssertExpectations(t) == false {
		t.Fatalf("Expectations not met")
	}

}

func TestSendResultError(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	sm := p.(*mockProcess).On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&sqs.SendMessageOutput{}, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
		sm.Unset()
	}()
	SendResult(types.QueueRequest{}, types.ImageResponseWrapper{})
	if p.(*mockProcess).AssertExpectations(t) == false {
		t.Fatalf("Expectations not met")
	}

}

func TestSendRequestToQueue(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	sm := p.(*mockProcess).On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&sqs.SendMessageOutput{}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
		sm.Unset()
	}()
	SendRequestToQueue(types.QueueRequest{})
	if p.(*mockProcess).AssertExpectations(t) == false {
		t.Fatalf("Expectations not met")
	}

}

func TestSendRequestToQueueError(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	sm := p.(*mockProcess).On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&sqs.SendMessageOutput{}, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
		sm.Unset()
	}()
	SendRequestToQueue(types.QueueRequest{})
	if p.(*mockProcess).AssertExpectations(t) == false {
		t.Fatalf("Expectations not met")
	}

}

func TestGetTraceId(t *testing.T) {
	traceId := GetTraceId()
	if traceId == "" {
		t.Fatalf("TraceId should not be empty")
	}
	if len(traceId) != 16 {
		t.Fatalf("TraceId should be 16 characters long")
	}
	if strings.Split(traceId, "-")[0] != traceId {
		t.Fatalf("TraceId should not contain dashes")
	}
}
func TestSubmitXRayTraceSubSegment(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	sm := p.(*mockProcess).On("PutTraceSegments", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&xray.PutTraceSegmentsOutput{}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
		sm.Unset()
	}()
	SubmitXRayTraceSubSegment("", "")
}

func TestSubmitXRayTraceSubSegmentError(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	sm := p.(*mockProcess).On("PutTraceSegments", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&xray.PutTraceSegmentsOutput{}, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
		sm.Unset()
	}()
	SubmitXRayTraceSubSegment("", "")
}

func TestStoreAnalyticsItem(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	pi := p.(*mockProcess).On("PutItem", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
		pi.Unset()
	}()
	StoreAnalyticsItem(types.AnalyticsItem{})
	if p.(*mockProcess).AssertExpectations(t) == false {
		t.Fatalf("Expectations not met")
	}

}

func TestStoreAnalyticsItemError(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	pi := p.(*mockProcess).On("PutItem", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
		pi.Unset()
	}()
	StoreAnalyticsItem(types.AnalyticsItem{})
	if p.(*mockProcess).AssertExpectations(t) == false {
		t.Fatalf("Expectations not met")
	}

}

func TestGetSecretValue(t *testing.T) {
	testString := "test"
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	gsv := p.(*mockProcess).On("GetSecretValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&secretsmanager.GetSecretValueOutput{SecretString: &testString}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
		gsv.Unset()
	}()
	GetSecretValue("")

}

func TestGetSecretValueError(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return errors.New("test error")
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	gsv := p.(*mockProcess).On("GetSecretValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
		gsv.Unset()
	}()
	GetSecretValue("")

}

func TestGetImageFromS3(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return nil
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	gis := p.(*mockProcess).On("GetS3Object", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{}, nil)
	defer func() {
		wr.Unset()
		ldc.Unset()
		gis.Unset()
	}()
	GetImageFromS3("")

}

func TestGetImageFromS3Error(t *testing.T) {
	wr := p.(*mockProcess).On("WithRegion", mock.Anything).Return(config.LoadOptionsFunc(func(o *config.LoadOptions) error {
		return errors.New("test error")
	}))
	ldc := p.(*mockProcess).On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{Region: region}, nil)
	gis := p.(*mockProcess).On("GetS3Object", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		wr.Unset()
		ldc.Unset()
		gis.Unset()
	}()
	GetImageFromS3("")

}

func TestSaveFile(t *testing.T) {
	testFile := &os.File{}
	nf := p.(*mockProcess).On("NewFile", mock.Anything, mock.Anything).Return(testFile)
	c := p.(*mockProcess).On("Copy", mock.Anything, mock.Anything).Return(int64(0), nil)
	defer func() {
		nf.Unset()
		c.Unset()
	}()
	SaveFile("", testFile)

}

func TestSaveFileError(t *testing.T) {
	testFile := &os.File{}
	nf := p.(*mockProcess).On("NewFile", mock.Anything, mock.Anything).Return(testFile)
	c := p.(*mockProcess).On("Copy", mock.Anything, mock.Anything).Return(int64(0), errors.New("test error"))
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Should have panicked")
		}
		nf.Unset()
		c.Unset()
	}()
	SaveFile("", testFile)
}
