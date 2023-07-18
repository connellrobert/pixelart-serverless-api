output "analytics_table_name" {
    value = aws_dynamodb_table.analytics_table.name
}

output "analytics_table_arn" {
    value = aws_dynamodb_table.analytics_table.arn
}

output "api_gateway_id" {
    value = aws_api_gateway_rest_api.pixelart_api.id
}

output "result_queue_name" {
    value = aws_sqs_queue.results_queue.name
}

output "result_queue_url" {
    value = aws_sqs_queue.results_queue.url
}

output "result_queue_arn" {
    value = aws_sqs_queue.results_queue.arn
}

# output "openai_secret_name" {
#     value = aws_secretsmanager_secret.openai_secret.name
# }

# output "openai_secret_arn" {
#     value = aws_secretsmanager_secret.openai_secret.arn
# }

output "api_gateway_url" {
    value = aws_api_gateway_deployment.pixelart_api_deployment.invoke_url
}