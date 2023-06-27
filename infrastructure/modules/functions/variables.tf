variable "lambda_source_path" {
    type = string
    default = "../"
}

variable "lambda_name" {
    type = string
    default = ""
}

variable "lambda_type" {
    type = string
    default = ""
}

variable "deployment_bucket_name" {
    type = string
    default = ""
}

variable "lambda_environment" {
    type = map(string)
}

