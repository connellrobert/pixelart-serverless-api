package ai

import (
	"io"
	"os"
	"testing"

	"github.com/aimless-it/ai-canvas/functions/lib/types"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
)

type mockProcessor struct {
	mock.Mock
}

func (m mockProcessor) GenerateImage(client openai.Client, request openai.ImageRequest) (openai.ImageResponse, error) {
	args := m.Called(client, request)
	return args.Get(0).(openai.ImageResponse), args.Error(1)
}

func (m mockProcessor) EditImage(client openai.Client, request openai.ImageEditRequest) (openai.ImageResponse, error) {
	args := m.Called(client, request)
	return args.Get(0).(openai.ImageResponse), args.Error(1)
}

func (m mockProcessor) CreateImageVariation(client openai.Client, request openai.ImageVariRequest) (openai.ImageResponse, error) {
	args := m.Called(client, request)
	return args.Get(0).(openai.ImageResponse), args.Error(1)
}

func (m mockProcessor) GetSecretValue(secretId string) string {
	args := m.Called(secretId)
	return args.String(0)
}

func (m mockProcessor) GetImageFromS3(image string) io.ReadCloser {
	args := m.Called(image)
	return args.Get(0).(*os.File)
}

func (m mockProcessor) SaveFile(fileName string, fileContents io.ReadCloser) *os.File {
	args := m.Called(fileName, fileContents)
	return args.Get(0).(*os.File)
}

func init() {
	ai = &mockProcessor{}
}

func TestGenerateImage(t *testing.T) {
	gsv := ai.(*mockProcessor).On("GetSecretValue", mock.Anything).Return("test-secret")
	gi := ai.(*mockProcessor).On("GenerateImage", mock.Anything, mock.Anything).Return(openai.ImageResponse{}, nil)
	defer func() {
		gsv.Unset()
		gi.Unset()
	}()
	giInput := types.GenerateImageRequest{
		Prompt:         "This is a test prompt",
		N:              1,
		Size:           "512x512",
		ResponseFormat: "image",
		User:           "test-user",
	}

	GenerateImage(giInput)
	ai.(*mockProcessor).AssertExpectations(t)
}

func TestEditImage(t *testing.T) {
	gsv := ai.(*mockProcessor).On("GetSecretValue", mock.Anything).Return("test-secret")
	gis := ai.(*mockProcessor).On("GetImageFromS3", mock.Anything).Return(&os.File{})
	sf := ai.(*mockProcessor).On("SaveFile", mock.Anything, mock.Anything).Return(&os.File{})
	ei := ai.(*mockProcessor).On("EditImage", mock.Anything, mock.Anything).Return(openai.ImageResponse{}, nil)
	defer func() {
		gsv.Unset()
		gis.Unset()
		ei.Unset()
		sf.Unset()
	}()

	eiInput := types.EditImageRequest{
		Image:          "test-image",
		Prompt:         "This is a test prompt",
		N:              1,
		Size:           "512x512",
		ResponseFormat: "image",
		Mask:           "test-mask",
		User:           "test-user",
	}
	EditImage(eiInput)
	ai.(*mockProcessor).AssertExpectations(t)
}

func TestCreateImageVariation(t *testing.T) {
	gsv := ai.(*mockProcessor).On("GetSecretValue", mock.Anything).Return("test-secret")
	gis := ai.(*mockProcessor).On("GetImageFromS3", mock.Anything).Return(&os.File{})
	sf := ai.(*mockProcessor).On("SaveFile", mock.Anything, mock.Anything).Return(&os.File{})
	civ := ai.(*mockProcessor).On("CreateImageVariation", mock.Anything, mock.Anything).Return(openai.ImageResponse{}, nil)
	defer func() {
		gsv.Unset()
		gis.Unset()
		civ.Unset()
		sf.Unset()
	}()

	civInput := types.CreateImageVariantRequest{
		Image:          "test-image",
		N:              1,
		Size:           "512x512",
		ResponseFormat: "image",
		User:           "test-user",
	}
	CreateImageVariation(civInput)
	ai.(*mockProcessor).AssertExpectations(t)
}
