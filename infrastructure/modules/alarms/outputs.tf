output "db_low_count_alarm_arn" {
  value = aws_cloudwatch_composite_alarm.db_low_count_alarm.arn
}

output "db_low_count_alarm_name" {
  value = aws_cloudwatch_composite_alarm.db_low_count_alarm.id
}

output "sns_topic_arn" {
    value = aws_sns_topic.sns_topic.arn
}

output "low_sqs_message_count_alarm_name" {
  value = aws_cloudwatch_metric_alarm.low_sqs_message_count_alarm.id
}

output "low_sqs_message_count_alarm_arn" {
  value = aws_cloudwatch_metric_alarm.low_sqs_message_count_alarm.arn
}