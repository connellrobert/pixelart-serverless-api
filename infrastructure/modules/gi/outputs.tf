output "gi_table_name" {
    value = local.table_name
}

output "gi_table_arn" {
    value = module.gi_queueing_system.queue_table_arn
}

output "gi_empty_db_alarm_arn" {
    value = module.gi_function_alarms.db_low_count_alarm_arn
}