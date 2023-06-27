
// Create a subscription to the sns topic
resource "aws_sns_topic_subscription" "sns_topic_subscription" {
  topic_arn = var.sns_topic_arn
  protocol  = "lambda"
  endpoint  = var.poll_function_arn
}


resource "aws_lambda_event_source_mapping" "oracle_trigger" {
  event_source_arn = var.oracle_queue_arn
  function_name = var.oracle_function_arn
  batch_size = 1
}