resource "aws_secretsmanager_secret_version" "openai_secret" {
  secret_id = aws_secretsmanager_secret.openai_secret.id
  secret_string = var.openai_api_key
}

resource "aws_secretsmanager_secret" "openai_secret" {
  name = "openai_api_token"
}