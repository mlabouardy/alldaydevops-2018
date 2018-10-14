resource "aws_iam_role" "iam_for_lambda_create" {
  name = "CreateInfrastructureRole"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "iam_policy_create" {
  name   = "CreateInfrastructurePolicy"
  policy = "${file("../create/policy.json")}"
}

resource "aws_iam_policy_attachment" "iam_policy_attachment_create" {
  name       = "attachment-policy-create"
  roles      = ["${aws_iam_role.iam_for_lambda_create.name}"]
  policy_arn = "${aws_iam_policy.iam_policy_create.arn}"
}

resource "aws_iam_role" "iam_for_lambda_provision" {
  name = "ProvisionInfrastructureRole"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "iam_policy_provision" {
  name   = "ProvisionInfrastructurePolicy"
  policy = "${file("../provision/policy.json")}"
}

resource "aws_iam_policy_attachment" "iam_policy_attachment_provision" {
  name       = "attachment-policy-provision"
  roles      = ["${aws_iam_role.iam_for_lambda_provision.name}"]
  policy_arn = "${aws_iam_policy.iam_policy_provision.arn}"
}

resource "aws_iam_instance_profile" "iam_for_ec2" {
  name = "SwarmClusterSSMRole"
  role = "${aws_iam_role.iam_role_for_ec2.name}"
}

resource "aws_iam_role" "iam_role_for_ec2" {
  name = "SwarmClusterSSMRole"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

data "aws_iam_policy" "AmazonEC2RoleforSSM" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_policy_attachment" "iam_policy_attachment_ec2" {
  name       = "attachment-policy-ec2"
  roles      = ["${aws_iam_role.iam_role_for_ec2.name}"]
  policy_arn = "${data.aws_iam_policy.AmazonEC2RoleforSSM.arn}"
}
