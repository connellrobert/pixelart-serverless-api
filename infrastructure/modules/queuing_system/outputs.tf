output "queue_url" {
    value = aws_sqs_queue.lambda_queue.url
}

# output "queue_table_arn" {
#     value = aws_dynamodb_table.queue_table.arn
# }

output "queue_arn" {
    value = aws_sqs_queue.lambda_queue.arn
}