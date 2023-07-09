
module "status_function" {
    source = "../functions"
    lambda_name = "status"
    lambda_type = "status"
    lambda_source_path = "${var.lambda_source_path}/status/bin"
    deployment_bucket_name = var.deployment_bucket_name
    lambda_environment = {
        ANALYTICS_TABLE_NAME = aws_dynamodb_table.analytics_table.name
    }
}


module "scheduler_function" {
    source = "../functions"
    lambda_name = "scheduler"
    lambda_type = "scheduler"
    lambda_source_path = "${var.lambda_source_path}/scheduler/bin"
    deployment_bucket_name = var.deployment_bucket_name
    lambda_environment = {
        "ANALYTICS_TABLE_NAME" = aws_dynamodb_table.analytics_table.name
        "GI_TABLE_NAME" = "gi_table"
    }
}
