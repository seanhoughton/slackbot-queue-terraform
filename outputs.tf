output event_url {
  value = "${aws_api_gateway_deployment.prod.invoke_url}/event"
}
