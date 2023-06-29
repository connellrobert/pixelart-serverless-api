// Create a cloudwatch alarm 
resource "aws_cloudwatch_metric_alarm" "low_sqs_message_count_alarm" {
  alarm_name          = local.low_sqs_message_count_alarm_name
  comparison_operator = "LessThanOrEqualToThreshold"
  metric_name         = "ApproximateNumberOfMessagesVisible"
  namespace           = "AWS/SQS"
  period              = "60"
  statistic           = "Sum"
  threshold           = "0"
  evaluation_periods  = "1"
  alarm_description   = "Trigger if the number of messages in the queue is 0"
  dimensions = {
    QueueName = var.queue_name
  }
}

// Create a composite cloudwatch alarm
resource "aws_cloudwatch_composite_alarm" "db_low_count_alarm" {
  alarm_name          = local.db_low_count_alarm_name
  alarm_description   = "Manually set to ALARM if the dynamodb count is 0 and to OK if records are added"
  alarm_rule          = "TRUE"
  depends_on = [  ]
}

// Create a composite alarm of the previous two alarms
resource "aws_cloudwatch_composite_alarm" "fill_queue_alarm" {
  alarm_name        = local.fill_queue_alarm_name
  alarm_description = "Triggers if GILowSQSMessageCountAlarm is triggered and GILowDynamoDBCountAlarm is OK"
  alarm_rule        = "ALARM(${aws_cloudwatch_metric_alarm.low_sqs_message_count_alarm.arn}) AND OK(${aws_cloudwatch_composite_alarm.db_low_count_alarm.arn})"
  alarm_actions     = [aws_sns_topic.sns_topic.arn]
  depends_on = [
    aws_cloudwatch_composite_alarm.db_low_count_alarm,
    aws_cloudwatch_metric_alarm.low_sqs_message_count_alarm
  ]
}

// Create a sns topic
resource "aws_sns_topic" "sns_topic" {
  name = var.sns_topic_name
}
locals {
  low_sqs_message_count_alarm_name = "${var.alarm_prefix}LowSQSMessageCountAlarm"
  db_low_count_alarm_name = "${var.alarm_prefix}LowDynamoDBCountAlarm"
  fill_queue_alarm_name = "${var.alarm_prefix}FillQueueAlarm"
}
