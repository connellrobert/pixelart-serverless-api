module "gi_function_crew" {
    source = "./modules/gi"
    lambda_source_path = var.lambda_source_path
    deployment_bucket_name = aws_s3_bucket.deployment_bucket.id
    openai_key_name = var.openai_key_name
    result_queue_url = module.core.result_queue_url
}

module "result_function" {
    source = "./modules/functions"
    lambda_name = "result"
    lambda_type = "result"
    lambda_source_path = "${var.lambda_source_path}/result/bin"
    deployment_bucket_name = aws_s3_bucket.deployment_bucket.id
    lambda_environment = {
        GENERATE_IMAGE_TABLE_NAME = module.gi_function_crew.gi_table_name
    }
}

resource "aws_iam_policy" "result_iam_policy" {
    name = "result_function_policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
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
            "Sid": "AllowDynamoDB",
            "Effect": "Allow",
            "Action": [
                "dynamodb:PutItem"
            ],
            "Resource": [
                "${module.gi_function_crew.gi_table_arn}"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "result_role_policy" {
    role = module.result_function.function_iam_role_name
    policy_arn = aws_iam_policy.result_iam_policy.arn
}
module "core" {
    source = "./modules/core"
}


resource "aws_lambda_event_source_mapping" "result_trigger" {
  event_source_arn = module.core.result_queue_arn
  function_name = module.result_function.function_arn
  batch_size = 10
}