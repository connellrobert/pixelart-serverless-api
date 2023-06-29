package lib

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	openai "github.com/sashabaranov/go-openai"
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
	Metadata             CommonMetadata
	Id                   string
	Action               RequestAction
	Priority             int
	CreateImage          GenerateImageRequest
	CreateImageEdit      EditImageRequest
	CreateImageVariation CreateImageVariantRequest
}

type ImageResponseWrapper struct {
	Response openai.ImageResponse `json:"response"`
	Success  bool                 `json:"success"`
}

type ResultRequest struct {
	Record QueueRequest         `json:"record"`
	Result ImageResponseWrapper `json:"result"`
}

// Analytics Item
type AnalyticsItem struct {
	Success  bool                            `json:"success"`
	Id       string                          `json:"id"`
	Record   QueueRequest                    `json:"record"`
	Attempts map[string]ImageResponseWrapper `json:"attempts"`
}

// Create dynamodb mappings for AnalyticsItem
func (r *AnalyticsItem) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: r.Id,
		},
		"request": &types.AttributeValueMemberM{
			Value: r.Record.ToDynamoDB(),
		},
		"attempts": &types.AttributeValueMemberM{
			Value: r.AttemptsToDynamoDB(),
		},
	}
}

func (r *AnalyticsItem) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Id = item["id"].(*types.AttributeValueMemberS).Value
	request := item["Record"].(*types.AttributeValueMemberM).Value
	record := request["request"].(*types.AttributeValueMemberM).Value
	action, err := strconv.Atoi(record["Action"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		panic(err)
	}
	r.Record.Action = RequestAction(action)
	switch r.Record.Action {
	case GenerateImageAction:
		r.Record.CreateImage.FromDynamoDB(record["CreateImage"].(*types.AttributeValueMemberM).Value)
	case EditImageAction:
		r.Record.CreateImageEdit.FromDynamoDB(record["CreateImageEdit"].(*types.AttributeValueMemberM).Value)
	case VariateImageAction:
		r.Record.CreateImageVariation.FromDynamoDB(record["CreateImageVariation"].(*types.AttributeValueMemberM).Value)
	}
	r.AttemptsFromDynamoDB(request["attempts"].(*types.AttributeValueMemberM).Value)
}

func (r *AnalyticsItem) AttemptsToDynamoDB() map[string]types.AttributeValue {
	attempts := make(map[string]types.AttributeValue)
	for k, v := range r.Attempts {
		attempts[k] = &types.AttributeValueMemberM{
			Value: v.ToDynamoDB(),
		}
	}
	return attempts
}

func (r *AnalyticsItem) AttemptsFromDynamoDB(item map[string]types.AttributeValue) {
	r.Attempts = make(map[string]ImageResponseWrapper)
	for k, v := range item {
		var irw ImageResponseWrapper
		irw.FromDynamoDB(v.(*types.AttributeValueMemberM).Value)
		r.Attempts[k] = irw
	}
}

// Create dynamodb mappings for types
func (r *QueueRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: r.Id,
		},
		"priority": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.Priority),
		},
		"Action": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(r.Action)),
		},
		"CreateImage": &types.AttributeValueMemberM{
			Value: r.CreateImage.ToDynamoDB(),
		},
		"CreateImageEdit": &types.AttributeValueMemberM{
			Value: r.CreateImageEdit.ToDynamoDB(),
		},
		"CreateImageVariation": &types.AttributeValueMemberM{
			Value: r.CreateImageVariation.ToDynamoDB(),
		},
	}
}

func (r *QueueRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	if item["id"] == nil {
		r.Id = item["PK"].(*types.AttributeValueMemberS).Value
	} else {
		r.Id = item["id"].(*types.AttributeValueMemberS).Value
	}
	r.Priority, _ = strconv.Atoi(item["priority"].(*types.AttributeValueMemberN).Value)
	action, err := strconv.Atoi(item["request"].(*types.AttributeValueMemberM).Value["Action"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		panic(err)
	}
	r.Action = RequestAction(action)
	switch r.Action {
	case GenerateImageAction:
		r.CreateImage.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["CreateImage"].(*types.AttributeValueMemberM).Value)
	case EditImageAction:
		r.CreateImageEdit.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["CreateImageEdit"].(*types.AttributeValueMemberM).Value)
	case VariateImageAction:
		r.CreateImageVariation.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["CreateImageVariation"].(*types.AttributeValueMemberM).Value)
	}
}

func (r *ResultRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: r.Record.Id,
		},
		"Priority": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.Record.Priority),
		},
		"request": &types.AttributeValueMemberM{
			Value: r.Record.ToDynamoDB(),
		},
		"result": &types.AttributeValueMemberM{
			Value: r.Result.ToDynamoDB(),
		},
	}
}

func (r *ResultRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Record.Id = item["PK"].(*types.AttributeValueMemberS).Value
	r.Record.Priority, _ = strconv.Atoi(item["Priority"].(*types.AttributeValueMemberN).Value)
	action, err := strconv.Atoi(item["request"].(*types.AttributeValueMemberM).Value["Action"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		panic(err)
	}

	r.Record.Action = RequestAction(action)
	switch r.Record.Action {
	case GenerateImageAction:
		r.Record.CreateImage.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["CreateImage"].(*types.AttributeValueMemberM).Value)
	case EditImageAction:
		r.Record.CreateImageEdit.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["CreateImageEdit"].(*types.AttributeValueMemberM).Value)
	case VariateImageAction:
		r.Record.CreateImageVariation.FromDynamoDB(item["request"].(*types.AttributeValueMemberM).Value["CreateImageVariation"].(*types.AttributeValueMemberM).Value)
	}
	r.Result.FromDynamoDB(item["result"].(*types.AttributeValueMemberM).Value)
}

func (r *GenerateImageRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Prompt": &types.AttributeValueMemberS{
			Value: r.Prompt,
		},
		"N": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.N),
		},
		"Size": &types.AttributeValueMemberS{
			Value: string(r.Size),
		},
		"ResponseFormat": &types.AttributeValueMemberS{
			Value: string(r.ResponseFormat),
		},
		"User": &types.AttributeValueMemberS{
			Value: r.User,
		},
	}
}

func (r *GenerateImageRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Prompt = item["Prompt"].(*types.AttributeValueMemberS).Value
	r.N, _ = strconv.Atoi(item["N"].(*types.AttributeValueMemberN).Value)
	r.Size = ImageSize(item["Size"].(*types.AttributeValueMemberS).Value)
	r.ResponseFormat = ResponseFormat(item["ResponseFormat"].(*types.AttributeValueMemberS).Value)
	r.User = item["User"].(*types.AttributeValueMemberS).Value
}

func (r *EditImageRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Prompt": &types.AttributeValueMemberS{
			Value: r.Prompt,
		},
		"N": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.N),
		},
		"Size": &types.AttributeValueMemberS{
			Value: string(r.Size),
		},
		"ResponseFormat": &types.AttributeValueMemberS{
			Value: string(r.ResponseFormat),
		},
		"User": &types.AttributeValueMemberS{
			Value: r.User,
		},
		"Image": &types.AttributeValueMemberS{
			Value: r.Image,
		},
		"Mask": &types.AttributeValueMemberS{
			Value: r.Mask,
		},
	}
}

func (r *EditImageRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.Prompt = item["Prompt"].(*types.AttributeValueMemberS).Value
	r.N, _ = strconv.Atoi(item["N"].(*types.AttributeValueMemberN).Value)
	r.Size = ImageSize(item["Size"].(*types.AttributeValueMemberS).Value)
	r.ResponseFormat = ResponseFormat(item["ResponseFormat"].(*types.AttributeValueMemberS).Value)
	r.User = item["User"].(*types.AttributeValueMemberS).Value
	r.Image = item["Image"].(*types.AttributeValueMemberS).Value
	r.Mask = item["Mask"].(*types.AttributeValueMemberS).Value
}

func (r *CreateImageVariantRequest) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"N": &types.AttributeValueMemberN{
			Value: strconv.Itoa(r.N),
		},
		"Size": &types.AttributeValueMemberS{
			Value: string(r.Size),
		},
		"ResponseFormat": &types.AttributeValueMemberS{
			Value: string(r.ResponseFormat),
		},
		"User": &types.AttributeValueMemberS{
			Value: r.User,
		},
		"Image": &types.AttributeValueMemberS{
			Value: r.Image,
		},
	}
}

func (r *CreateImageVariantRequest) FromDynamoDB(item map[string]types.AttributeValue) {
	r.N, _ = strconv.Atoi(item["N"].(*types.AttributeValueMemberN).Value)
	r.Size = ImageSize(item["Size"].(*types.AttributeValueMemberS).Value)
	r.ResponseFormat = ResponseFormat(item["ResponseFormat"].(*types.AttributeValueMemberS).Value)
	r.User = item["User"].(*types.AttributeValueMemberS).Value
	r.Image = item["Image"].(*types.AttributeValueMemberS).Value
}

// Create dynamodb AttributeValue List for openai response
func (r *ImageResponseWrapper) MapDataInnerToDynamoDB() []string {
	lst := make([]string, 0)
	for _, d := range r.Response.Data {
		data, err := json.Marshal(d)
		if err != nil {
			panic(err)
		}
		lst = append(lst, string(data))
	}
	return lst
}

func (r *ImageResponseWrapper) ToDynamoDB() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Created": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(r.Response.Created)),
		},
		"Data": &types.AttributeValueMemberSS{
			Value: r.MapDataInnerToDynamoDB(),
		},
	}
}

func (r *ImageResponseWrapper) FromDynamoDB(item map[string]types.AttributeValue) {
	tmp, _ := strconv.Atoi(item["Created"].(*types.AttributeValueMemberN).Value)
	r.Response.Created = int64(tmp)
	r.Response.Data = make([]openai.ImageResponseDataInner, 0)
	for _, d := range item["Data"].(*types.AttributeValueMemberSS).Value {
		var data openai.ImageResponseDataInner
		err := json.Unmarshal([]byte(d), &data)
		if err != nil {
			panic(err)
		}
		r.Response.Data = append(r.Response.Data, data)
	}
}

func (q *QueueRequest) MapParams(action RequestAction, params interface{}) {
	switch action {
	case GenerateImageAction:
		q.CreateImage = GenerateImageRequest{
			Prompt:         params.(map[string]interface{})["prompt"].(string),
			N:              int(params.(map[string]interface{})["n"].(float64)),
			Size:           ImageSize(params.(map[string]interface{})["size"].(string)),
			ResponseFormat: ResponseFormat(params.(map[string]interface{})["responseFormat"].(string)),
			User:           params.(map[string]interface{})["user"].(string),
		}
	case EditImageAction:
		q.CreateImageEdit = EditImageRequest{
			Prompt:         params.(map[string]interface{})["prompt"].(string),
			N:              int(params.(map[string]interface{})["n"].(float64)),
			Size:           ImageSize(params.(map[string]interface{})["size"].(string)),
			ResponseFormat: ResponseFormat(params.(map[string]interface{})["responseFormat"].(string)),
			User:           params.(map[string]interface{})["user"].(string),
			Image:          params.(map[string]interface{})["image"].(string),
			Mask:           params.(map[string]interface{})["mask"].(string),
		}
	case VariateImageAction:
		q.CreateImageVariation = CreateImageVariantRequest{
			N:              int(params.(map[string]interface{})["n"].(float64)),
			Size:           ImageSize(params.(map[string]interface{})["size"].(string)),
			ResponseFormat: ResponseFormat(params.(map[string]interface{})["responseFormat"].(string)),
			User:           params.(map[string]interface{})["user"].(string),
			Image:          params.(map[string]interface{})["image"].(string),
		}
	}
}
