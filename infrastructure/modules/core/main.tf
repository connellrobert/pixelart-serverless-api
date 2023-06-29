// Create an api gateway instance
resource "aws_api_gateway_rest_api" "pixelart_api" {
    name = "pixelart-api"
}

resource "aws_api_gateway_resource" "image" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    parent_id   = aws_api_gateway_rest_api.pixelart_api.root_resource_id
    path_part   = "image"
}

resource "aws_api_gateway_method" "image_post" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    resource_id = aws_api_gateway_resource.image.id
    http_method = "POST"
    authorization = "NONE"
}


# resource "aws_api_gateway_deployment" "pixelart_api_deployment" {
#   rest_api_id = aws_api_gateway_rest_api.pixelart_api.id

#   triggers = {
#     redeployment = "testing for now, replace with something that changes in the api"
#   }

#   lifecycle {
#     create_before_destroy = true
#   }
# }

# resource "aws_api_gateway_stage" "api_conection" {
#   deployment_id = aws_api_gateway_deployment.pixelart_api_deployment.id
#   rest_api_id   = aws_api_gateway_rest_api.pixelart_api.id
#   stage_name    = "test"
# }

// Create dynamodb table
resource "aws_dynamodb_table" "analytics_table" {
    name           = "analytics_table"
    billing_mode   = "PAY_PER_REQUEST"
    hash_key       = "id"
    attribute {
        name = "id"
        type = "S"
    }
}

// Create a sqs queue for lambda function results
resource "aws_sqs_queue" "results_queue" {
    name = "results_queue"
}

module "scheduler_function" {
    source = "../functions"
    lambda_name = "scheduler"
    lambda_type = "scheduler"
    lambda_source_path = var.lambda_source_path
    deployment_bucket_name = var.deployment_bucket_name
    lambda_environment = {
        "ANALYTICS_TABLE_NAME" = aws_dynamodb_table.analytics_table.name
        "GI_TABLE_NAME" = "gi_table"
    }
}

resource "aws_api_gateway_integration" "sheduler_integration" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    resource_id = aws_api_gateway_resource.image.id
    http_method = aws_api_gateway_method.image_post.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = module.scheduler_function.function_invoke_arn
}

resource "aws_lambda_permission" "scheduler_permission" {
    statement_id = "AllowAPIGatewayInvoke"
    action = "lambda:InvokeFunction"
    function_name = module.scheduler_function.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.pixelart_api.execution_arn}/*/*/*"
}

resource "aws_iam_role_policy_attachment" "scheduler_role_policy" {
    role = module.scheduler_function.function_iam_role_name
    policy_arn = aws_iam_policy.scheduler_policy.arn
}

resource "aws_iam_policy" "scheduler_policy" {
    name = "scheduler_policy"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "dynamodb:PutItem"
            ],
            "Effect": "Allow",
            "Resource": [
                "${aws_dynamodb_table.analytics_table.arn}",
                "${var.gi_table_arn}"
            ]
        },
        {
            "Action": [
                "cloudwatch:SetAlarmState"
            ],
            "Effect": "Allow",
            "Resource": [
                "${var.gi_empty_db_alarm_arn}"
            ]
        },
        {
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Effect": "Allow",
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