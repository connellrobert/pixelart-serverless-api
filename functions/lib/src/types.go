package lib

import (
	openai "github.com/sashabaranov/go-openai"
)

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

func (rf ResponseFormat) openaiResponseFormat() string {
	switch rf {
	case URL:
		return openai.CreateImageResponseFormatURL
	case BASE64:
		return openai.CreateImageResponseFormatB64JSON
	default:
		return openai.CreateImageResponseFormatURL
	}
}

type GenerateImageRequest struct {
	Prompt         string
	N              int
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
}
type EditImageRequest struct {
	Prompt         string
	N              int
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
	Image          string
	Mask           string
}

type CreateImageVariantRequest struct {
	N              int
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
	Image          string
}

type RequestAction int

const (
	GenerateImageAction RequestAction = iota
	EditImageAction
	VariateImageAction
)

type QueueRequest struct {
	Id                   string
	Action               RequestAction
	Priority             int
	CreateImage          GenerateImageRequest
	CreateImageEdit      EditImageRequest
	CreateImageVariation CreateImageVariantRequest
}

type ResultRequest struct {
	Record QueueRequest         `json:"record"`
	Result openai.ImageResponse `json:"result"`
}
