
// Create a dynamodb table for the lambda function
# resource "aws_dynamodb_table" "queue_table" {
#   name           = var.table_name
#   billing_mode   = "PAY_PER_REQUEST"
#   hash_key       = "id"
#   range_key = "priority"
#   attribute {
#     name = "id"
#     type = "S"
#   }
#   attribute {
#     name = "priority"
#     type = "N"
#   }
# }

// Create a sqs queue for the lambda function
resource "aws_sqs_queue" "lambda_queue" {
  name = var.queue_name
  delay_seconds = 0
  max_message_size = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 0
  visibility_timeout_seconds = 66
  fifo_queue = true
}
