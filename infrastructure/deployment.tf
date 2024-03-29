// create a s3 resource
resource "random_string" "deployment_suffix" {
  length = 4
  special = false
  upper = false
  number = false
}

resource "random_string" "customer_data_suffix" {
  length = 4
  special = false
  upper = false
  number = false
}

resource "aws_s3_bucket" "deployment_bucket" {
  bucket = "${var.deployment_bucket_name}-${random_string.deployment_suffix.result}"
  force_destroy = true
}

resource "aws_s3_bucket" "customer_data_bucket" {
    bucket = "${var.customer-data-bucket-name}-${random_string.customer_data_suffix.result}"
    force_destroy = true
}

module "core" {
    source = "./modules/core"
    lambda_source_path = var.lambda_source_path
    deployment_bucket_name = aws_s3_bucket.deployment_bucket.id
    hosted_zone_id = var.route_53_hosted_zone_id
    queue_url = module.function_queue.queue_url
    queue_arn = module.function_queue.queue_arn
    route53_domain = var.route53_domain
}

module "function_queue" {
    source = "./modules/queuing_system"
    queue_name = "pa-queue.fifo"
}


module "result_function" {
    source = "./modules/functions"
    lambda_name = "result"
    lambda_type = "result"
    lambda_source_path = "${var.lambda_source_path}/result/bin"
    deployment_bucket_name = aws_s3_bucket.deployment_bucket.id
    lambda_environment = {
        ANALYTICS_TABLE_NAME = module.core.analytics_table_name
        QUEUE_URL = module.function_queue.queue_url
        IMAGE_BUCKET = aws_s3_bucket.customer_data_bucket.id
    }
}

resource "aws_iam_policy" "result_iam_policy" {
    name = "result_function_policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "dynamodb:GetItem",
                "dynamodb:UpdateItem"
            ],
            "Effect": "Allow",
            "Resource": [
                "${module.core.analytics_table_arn}"
            ]
        },
        {
            "Sid": "AllowSQS",
            "Effect": "Allow",
            "Action": [
                "sqs:DeleteMessage",
                "sqs:ReceiveMessage",
                "sqs:GetQueueUrl",
                "sqs:GetQueueAttributes",
                "sqs:ListQueues"
            ],
            "Resource": [
                "${module.core.result_queue_arn}"
            ]
        },
        {
            "Action" : [
                "sqs:SendMessage"
            ],
            "Effect": "Allow",
            "Resource" : [
                "${module.function_queue.queue_arn}"
            ]
        },
        {
            "Sid": "AllowLogs",
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "*"
            ]
        },
        {
            "Action": "xray:*",
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "result_role_policy" {
    role = module.result_function.function_iam_role_name
    policy_arn = aws_iam_policy.result_iam_policy.arn
}


resource "aws_lambda_event_source_mapping" "result_trigger" {
  event_source_arn = module.core.result_queue_arn
  function_name = module.result_function.function_arn
  batch_size = 10
}

module "oracle_function" {
  source                 = "./modules/functions"
  lambda_name            = "oracle"
  lambda_type            = "oracle"
  lambda_source_path     = "${var.lambda_source_path}/oracle/bin"
  deployment_bucket_name = aws_s3_bucket.deployment_bucket.id
  lambda_environment = {
    "OPENAI_API_KEY_SECRET_ID" = aws_secretsmanager_secret.openai_secret.id
    "RESULT_QUEUE_URL" = module.core.result_queue_url
    "DEBUG_MODE"      = var.debug_mode
    "IMAGE_BUCKET" = aws_s3_bucket.customer_data_bucket.id
  }
}


resource "aws_lambda_event_source_mapping" "oracle_trigger" {
  event_source_arn = module.function_queue.queue_arn
  function_name = module.oracle_function.function_name
  batch_size = 1
}


resource "aws_iam_role_policy_attachment" "oracle_role_policy" {
  role       = module.oracle_function.function_iam_role_name
  policy_arn = aws_iam_policy.oracle_policy.arn
}

// Create an IAM policy with basic lambda execution permissions
resource "aws_iam_policy" "oracle_policy" {
  name = "oracle_policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": "${aws_s3_bucket.customer_data_bucket.arn}/*"
    },
    {
      "Action": [
        "sqs:DeleteMessage",
        "sqs:ReceiveMessage",
        "sqs:GetQueueUrl",
        "sqs:GetQueueAttributes",
        "sqs:ListQueues"
      ],
      "Effect": "Allow",
      "Resource": "${module.function_queue.queue_arn}"
    },
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "xray:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "sqs:SendMessage"
      ],
      "Effect": "Allow",
      "Resource": "${module.core.result_queue_arn}"
    },
    {
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Effect": "Allow",
      "Resource": "${aws_secretsmanager_secret.openai_secret.arn}"
    }
  ]
}
EOF
}
