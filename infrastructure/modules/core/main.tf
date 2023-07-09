
// Create dynamodb table
resource "aws_dynamodb_table" "analytics_table" {
    name           = "analytics_table"
    billing_mode   = "PAY_PER_REQUEST"
    hash_key       = "id"
    attribute {
        name = "id"
        type = "S"
    }
}

// Create a sqs queue for lambda function results
resource "aws_sqs_queue" "results_queue" {
    name = "results_queue"
}
