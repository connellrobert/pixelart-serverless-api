package types

type ImageSize string

const (
	Small  ImageSize = "256x256"
	Medium ImageSize = "512x512"
	Large  ImageSize = "1024x1024"
)

type ResponseFormat string

const (
	URL    ResponseFormat = "URL"
	BASE64 ResponseFormat = "BASE64"
)

type GenerateImageRequest struct {
	Prompt         string
	N              int
	Size           ImageSize
	ResponseFormat ResponseFormat
	User           string
}

type EditImageRequest struct {
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

type RequestAction string

const (
	GenerateImage RequestAction = "createImage"
	EditImage     RequestAction = "createImageEdit"
	VariateImage  RequestAction = "createImageVariation"
)

type QueueRequest struct {
	Id                   string
	Action               RequestAction
	Priority             int
	CreateImage          GenerateImageRequest
	CreateImageEdit      EditImageRequest
	CreateImageVariation CreateImageVariantRequest
}
