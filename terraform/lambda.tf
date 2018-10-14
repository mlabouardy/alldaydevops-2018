resource "aws_lambda_function" "lambda_create_infrastructure" {
  filename      = "../create/deployment.zip"
  function_name = "CreateInfrastructure"
  role          = "${aws_iam_role.iam_for_lambda_create.arn}"
  handler       = "create/main"
  runtime       = "go1.x"
  timeout       = 60

  environment {
    variables = {
      AMI               = "${var.docker_ami}"
      KEYPAIR           = "${var.keypair}"
      SECURITY_GROUP_ID = "${aws_security_group.security_group.id}"
      SQS_URL           = "${aws_sqs_queue.queue.id}"
      SSM_ROLE_NAME     = "${aws_iam_role.iam_role_for_ec2.name}"
      TABLE_NAME        = "${aws_dynamodb_table.table.name}"
    }
  }
}

resource "aws_lambda_function" "lambda_provision_infrastructure" {
  filename      = "../provision/deployment.zip"
  function_name = "ProvisionInfrastructure"
  role          = "${aws_iam_role.iam_for_lambda_provision.arn}"
  handler       = "provision/main"
  runtime       = "go1.x"
  timeout       = 120

  environment {
    variables = {
      SLACK_WEBHOOK = "${var.slack_webhook}"
      SQS_URL       = "${aws_sqs_queue.queue.id}"
      TABLE_NAME    = "${aws_dynamodb_table.table.name}"
    }
  }
}
