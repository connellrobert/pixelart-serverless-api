// Create an api gateway instance
resource "aws_api_gateway_rest_api" "pixelart_api" {
    name = "pixelart-api"
}

resource "aws_api_gateway_method_settings" "general_api_settings" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    stage_name = aws_api_gateway_stage.api_conection.stage_name
    method_path = "*/*"
    settings {
        logging_level = "INFO"
        metrics_enabled = true
        data_trace_enabled = true
        throttling_burst_limit = 5000
        throttling_rate_limit = 10000
    }
}

resource "aws_api_gateway_account" "api_gateway_account" {
    cloudwatch_role_arn = aws_iam_role.cloudwatch.arn

}


resource "aws_iam_role" "cloudwatch" {
  name = "pixelart-cloudwatch-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "cloudwatch" {
  name = "default"
  role = "${aws_iam_role.cloudwatch.id}"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams",
                "logs:PutLogEvents",
                "logs:GetLogEvents",
                "logs:FilterLogEvents"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "xray:*",
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_acm_certificate" "api_domain" {
    domain_name = "api.aimless.it"
    validation_method = "DNS"
    tags = {
        Name = "api.aimless.it"
    }
}

resource "aws_acm_certificate_validation" "api_domain_validation" {
    certificate_arn = aws_acm_certificate.api_domain.arn
    # validation_record_fqdns = [aws_route53_record.api.fqdn]
}

resource "aws_route53_record" "api" {
    name = aws_api_gateway_domain_name.api_domain.domain_name
    zone_id = var.hosted_zone_id
    type = "A"
    alias {
        name = aws_api_gateway_domain_name.api_domain.regional_domain_name
        zone_id = aws_api_gateway_domain_name.api_domain.regional_zone_id
        evaluate_target_health = false
    }
}

resource "aws_api_gateway_domain_name" "api_domain" {
    domain_name  = "api.aimless.it"
    regional_certificate_arn = aws_acm_certificate_validation.api_domain_validation.certificate_arn
    endpoint_configuration {
        types = ["REGIONAL"]
    }
}

resource "aws_api_gateway_base_path_mapping" "api_domain_mapping" {
    domain_name = aws_api_gateway_domain_name.api_domain.domain_name
    api_id = aws_api_gateway_rest_api.pixelart_api.id
    stage_name = aws_api_gateway_stage.api_conection.stage_name
}



// API resources
resource "aws_api_gateway_resource" "image" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    parent_id   = aws_api_gateway_rest_api.pixelart_api.root_resource_id
    path_part   = "image"
}

resource "aws_api_gateway_resource" "status" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    parent_id   = aws_api_gateway_rest_api.pixelart_api.root_resource_id
    path_part   = "status"
}

resource "aws_api_gateway_resource" "status_id" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    parent_id   = aws_api_gateway_resource.status.id
    path_part   = "{id}"
}

// API methods
resource "aws_api_gateway_method" "image_post" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    resource_id = aws_api_gateway_resource.image.id
    http_method = "POST"
    authorization = "NONE"
}

resource "aws_api_gateway_method" "status_get" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    resource_id = aws_api_gateway_resource.status_id.id
    http_method = "GET"
    authorization = "NONE"
}



resource "aws_api_gateway_deployment" "pixelart_api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.pixelart_api.id

  triggers = {
    redeployment = "testing for now, replace with something that changes in the api"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "api_conection" {
  deployment_id = aws_api_gateway_deployment.pixelart_api_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.pixelart_api.id
  stage_name    = "test"
  xray_tracing_enabled = true
}


// API integrations
resource "aws_api_gateway_integration" "sheduler_integration" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    resource_id = aws_api_gateway_resource.image.id
    http_method = aws_api_gateway_method.image_post.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = module.scheduler_function.function_invoke_arn
}

resource "aws_api_gateway_integration" "status_integration" {
    rest_api_id = aws_api_gateway_rest_api.pixelart_api.id
    resource_id = aws_api_gateway_resource.status_id.id
    http_method = aws_api_gateway_method.status_get.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = module.status_function.function_invoke_arn
}

// Invoke permissions


resource "aws_lambda_permission" "scheduler_permission" {
    statement_id = "AllowAPIGatewayInvoke"
    action = "lambda:InvokeFunction"
    function_name = module.scheduler_function.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.pixelart_api.execution_arn}/*/*/*"
}



resource "aws_lambda_permission" "status_permission" {
    statement_id = "AllowAPIGatewayInvoke"
    action = "lambda:InvokeFunction"
    function_name = module.status_function.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.pixelart_api.execution_arn}/*/*/*"
}
