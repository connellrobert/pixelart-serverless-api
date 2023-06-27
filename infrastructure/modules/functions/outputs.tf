output "function_arn" {
    value = aws_lambda_function.lambda_function.arn
}

output "function_name" {
    value = aws_lambda_function.lambda_function.function_name
}

output "function_invoke_arn" {
    value = aws_lambda_function.lambda_function.invoke_arn
}

output "function_iam_role_arn" {
    value = aws_iam_role.lambda_role.arn
}

output "function_iam_role_name" {
    value = aws_iam_role.lambda_role.name
}

