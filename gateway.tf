// Create an API Gateway with an `/event` endpoint
//

resource "aws_api_gateway_rest_api" "event_api" {
  name        = "${var.service_name}-api"
  description = "This provides a public api for Slack to post events into. Data is inserted into an SQS queue"
}

resource "aws_api_gateway_resource" "event" {
  rest_api_id = "${aws_api_gateway_rest_api.event_api.id}"
  parent_id   = "${aws_api_gateway_rest_api.event_api.root_resource_id}"
  path_part   = "event"
}

resource "aws_api_gateway_method" "event" {
  rest_api_id   = "${aws_api_gateway_rest_api.event_api.id}"
  resource_id   = "${aws_api_gateway_resource.event.id}"
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "event" {
  rest_api_id = "${aws_api_gateway_rest_api.event_api.id}"
  resource_id = "${aws_api_gateway_resource.event.id}"
  http_method = "${aws_api_gateway_method.event.http_method}"
  status_code = "200"
}

// Set up the /event endpoint to send data to SQS
//

resource "aws_api_gateway_integration" "gateway_to_sqs" {
  rest_api_id = "${aws_api_gateway_rest_api.event_api.id}"
  resource_id = "${aws_api_gateway_method.event.resource_id}"
  http_method = "${aws_api_gateway_method.event.http_method}"

  //credentials             = "${aws_iam_role.lambda_access.arn}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${aws_lambda_function.equeue_event.arn}/invocations"
}

// Let the gatway access the lambda
resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.equeue_event.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.event_api.id}/*/${aws_api_gateway_method.event.http_method}${aws_api_gateway_resource.event.path}"
}

// Create a deployment "prod"

resource "aws_api_gateway_deployment" "prod" {
  rest_api_id = "${aws_api_gateway_rest_api.event_api.id}"
  stage_name  = "prod"
}
