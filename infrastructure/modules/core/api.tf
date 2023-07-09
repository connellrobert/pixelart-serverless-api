// Create an api gateway instance
resource "aws_api_gateway_rest_api" "pixelart_api" {
    name = "pixelart-api"
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
