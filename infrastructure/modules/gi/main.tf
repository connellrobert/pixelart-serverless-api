
module "gi_oracle_function" {
  source                 = "../functions"
  lambda_name            = "gi_oracle"
  lambda_type            = "oracle"
  lambda_source_path     = "${var.lambda_source_path}/oracle/bin"
  deployment_bucket_name = var.deployment_bucket_name
  lambda_environment = {
    "OPENAI_API_KEY"   = var.openai_key_name
    "RESULT_QUEUE_URL" = var.result_queue_url
  }
}

module "gi_poll_function" {
  source                 = "../functions"
  lambda_name            = "gi_poll"
  lambda_type            = "poll"
  lambda_source_path     = "${var.lambda_source_path}/poll/bin"
  deployment_bucket_name = var.deployment_bucket_name
  lambda_environment = {
    "TABLE_NAME"          = local.table_name
    "QUEUE_URL"           = module.gi_queueing_system.queue_url
    "EMPTY_DB_ALARM_NAME" = module.gi_function_alarms.db_low_count_alarm_name
    "RESULT_QUEUE_URL"    = var.result_queue_url
  }
}

module "gi_function_policies" {
  source           = "../function_policies"
  oracle_role_name = module.gi_oracle_function.function_iam_role_name
  poll_role_name   = module.gi_poll_function.function_iam_role_name
  queue_arn        = module.gi_queueing_system.queue_arn
  table_arn        = module.gi_queueing_system.queue_table_arn
}

module "gi_queueing_system" {
  source     = "../queuing_system"
  queue_name = local.queue_name
  table_name = local.table_name
}

module "gi_function_alarms" {
  source         = "../alarms"
  sns_topic_name = "gi_alarm_topic"
  queue_name     = local.queue_name
  alarm_prefix   = "gi"
}



module "gi_triggers" {
  source              = "../function_triggers"
  poll_function_arn   = module.gi_poll_function.function_arn
  oracle_function_arn = module.gi_oracle_function.function_arn
  oracle_queue_arn    = module.gi_queueing_system.queue_arn
  sns_topic_arn       = module.gi_function_alarms.sns_topic_arn
}

locals {
  table_name = "gi_table"
  queue_name = "gi-queue.fifo"
}
