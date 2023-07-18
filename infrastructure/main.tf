terraform {
    required_version = ">= 0.12"
    backend "s3" {
        bucket = "terraform-shitty-shit-backend-stuff"
        key = "pixelart/tf-backend"
        region = "us-east-1"
    }
    required_providers {
        aws = {
        source  = "hashicorp/aws"
        version = "~> 3.0"
        }
    }
}

provider "aws" {
    region = "us-east-1"
}

variable "deployment_bucket_name" {
    type = string
    default = "my-bucket"
}

variable "customer-data-bucket-name" {
    type = string
    default = "my-customer-data-bucket"
}

# variable "functions" {
#     type = list(object({
#         name = string
#         type = string
#         source_path = string
#         environment = map(string)
#     }))
# }

variable "lambda_source_path" {
    type = string
    default = "/functions"
}


output "pixelart_api_url" {
    value = module.core.api_gateway_url
}

variable "debug_mode" {
    type = bool
    default = false
}

variable "route_53_hosted_zone_id" {
    type = string
}

variable "OPENAI_API_KEY" {
    type = string
}