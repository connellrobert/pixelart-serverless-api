data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "${var.lambda_source_path}"
  output_path = "${path.module}/${var.lambda_name}.zip"
}

locals {
  lambda_object_key = filemd5(data.archive_file.lambda_zip.output_path)
}

resource "aws_s3_bucket_object" "lambda_object" {
    bucket = var.deployment_bucket_name
    key    = "${local.lambda_object_key}.zip}"
    source = data.archive_file.lambda_zip.output_path
    etag   = local.lambda_object_key
}

resource "aws_lambda_function" "lambda_function" {
  function_name = var.lambda_name
  s3_bucket     = var.deployment_bucket_name
  s3_key        = aws_s3_bucket_object.lambda_object.id
  handler       = var.lambda_type
  runtime       = "go1.x"
  role          = aws_iam_role.lambda_role.arn
  environment {
    variables = var.lambda_environment
  }
}

resource "random_string" "lambda_role_suffix" {
  length  = 4
  special = false
}

resource "aws_iam_role" "lambda_role" {
  name = "lambda-${var.lambda_name}-${random_string.lambda_role_suffix.result}"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
