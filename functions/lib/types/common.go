package types

import openai "github.com/sashabaranov/go-openai"

type RequestAction int

const (
	GenerateImageAction RequestAction = iota
	EditImageAction
	VariateImageAction
)

type CommonMetadata struct {
	TraceId string `json:"tracing_id"`
}

type ImageSize string

const (
	Small  ImageSize = "256x256"
	Medium ImageSize = "512x512"
	Large  ImageSize = "1024x1024"
)

func (is ImageSize) OpenaiImageSize() string {
	switch is {
	case Small:
		return openai.CreateImageSize256x256
	case Medium:
		return openai.CreateImageSize512x512
	case Large:
		return openai.CreateImageSize1024x1024
	default:
		return openai.CreateImageSize256x256
	}
}

type ResponseFormat string

const (
	URL    ResponseFormat = "URL"
	BASE64 ResponseFormat = "BASE64"
)

func (rf ResponseFormat) OpenaiResponseFormat() string {
	switch rf {
	case URL:
		return openai.CreateImageResponseFormatURL
	case BASE64:
		return openai.CreateImageResponseFormatB64JSON
	default:
		return openai.CreateImageResponseFormatURL
	}
}
