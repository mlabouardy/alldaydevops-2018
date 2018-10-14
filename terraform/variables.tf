variable "region" {
  description = "AWS region"
  default     = "us-east-1"
}

variable "shared_credentials_file" {
  description = "AWS credentials file path"
  default     = "/Users/mlabouardy/.ssh/credentials"
}

variable "aws_profile" {
  description = "AWS profile"
  default     = "default"
}

variable "table_name" {
  description = "DynamoDB Table's name"
  default     = "clusters"
}

variable "queue_name" {
  description = "SQS name"
  default     = "clusters"
}

variable "keypair" {
  description = "SSH KeyPair"
}

variable "docker_ami" {
  description = "AMI with Docker engine preinstalled"
}

variable "slack_webhook" {
  description = "Slack Webhook URL"
}
