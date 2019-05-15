provider "archive" {}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "${path.module}/enqueue_event.py"
  output_path = "${path.root}/files/enqueue_event.zip"
}

data "aws_iam_policy_document" "lambda_policy" {
  statement {
    sid    = ""
    effect = "Allow"

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda_access" {
  name               = "${var.service_name}-lamda-access"
  assume_role_policy = "${data.aws_iam_policy_document.lambda_policy.json}"
}

resource "aws_iam_role_policy_attachment" "event_queue_role_attach" {
  role       = "${aws_iam_role.lambda_access.name}"
  policy_arn = "${aws_iam_policy.event_queue_policy.arn}"
}

resource "aws_lambda_function" "equeue_event" {
  function_name = "enqueue_${var.service_name}_event"

  filename         = "${data.archive_file.zip.output_path}"
  source_code_hash = "${data.archive_file.zip.output_sha}"

  role    = "${aws_iam_role.lambda_access.arn}"
  handler = "enqueue_event.lambda_handler"
  runtime = "python3.6"

  environment {
    variables = {
      queue_name = "${var.service_name}-events.fifo"
    }
  }
}
