# variable "gi_table_arn" {
#     type = string
# }

# variable "gi_empty_db_alarm_arn" {
#     type = string
# }

variable "lambda_source_path" {
    type = string
}

variable "deployment_bucket_name" {
    type = string
}

variable "hosted_zone_id" {
    type = string
}

variable "queue_url" {
    type = string
}

variable "queue_arn" {
    type = string
}