variable "deployment_bucket_name" {
    type = string
    default = "my-bucket"
}

variable "customer-data-bucket-name" {
    type = string
    default = "my-customer-data-bucket"
}

variable "lambda_source_path" {
    type = string
    default = "/functions"
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

variable "route53_domain" {
    type = string
}