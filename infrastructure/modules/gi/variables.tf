variable "lambda_source_path" {
    type = string
    default = "../"
}

variable "deployment_bucket_name" {
    type = string
    default = ""
}

variable "openai_key_name" {
    type = string
    default = ""
}


variable "result_queue_url" {
    type = string
    default = ""
}

variable "result_queue_arn" {
    type = string
}

variable "openai_secret_name" {
    type = string
}

variable "openai_secret_arn" {
    type = string
}

variable debug {
    type = bool
    default = false
}