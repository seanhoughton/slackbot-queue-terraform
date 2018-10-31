// Create an SQS queue for posted events
//

data "aws_iam_policy_document" "event_queue_policy" {
  statement {
    effect    = "Allow"
    resources = ["${aws_sqs_queue.event_queue.arn}"]

    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:GetQueueUrl",
    ]
  }
}

resource "aws_iam_policy" "event_queue_policy" {
  name        = "${var.service_name}-access-policy"
  path        = "/"
  description = "Read/write access to the ${var.service_name} queue"
  policy      = "${data.aws_iam_policy_document.event_queue_policy.json}"
}

resource "aws_sqs_queue" "event_queue" {
  name                        = "${var.service_name}-events.fifo"
  fifo_queue                  = true
  delay_seconds               = 1
  max_message_size            = 65536                             // slack maximum message size is 40,000 bytes
  message_retention_seconds   = 86400
  receive_wait_time_seconds   = 10
  content_based_deduplication = true

  tags {
    Environment = "production"
    Terraform   = "true"
    Service     = "${var.service_name}"
  }
}
