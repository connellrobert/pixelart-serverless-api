resource "aws_secretsmanager_secret_version" "openai_secret" {
  secret_id = aws_secretsmanager_secret.openai_secret.id
  secret_string = var.OPENAI_API_KEY
}

# create random string for secret name
resource "random_string" "openai_secret_name" {
  length = 4
  special = false
  upper = false
  number = false
}

resource "aws_secretsmanager_secret" "openai_secret" {
    name = "openai-secret-${random_string.openai_secret_name.result}"
}
