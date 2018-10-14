resource "aws_sqs_queue" "queue" {
  name                       = "${var.queue_name}"
  delay_seconds              = 90
  visibility_timeout_seconds = 130

  tags {
    Environment = "AllDayDevOps"
  }
}
