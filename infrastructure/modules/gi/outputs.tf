output "gi_table_name" {
    value = local.table_name
}

output "gi_table_arn" {
    value = module.gi_queueing_system.queue_table_arn
}