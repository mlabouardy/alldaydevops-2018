resource "aws_dynamodb_table" "table" {
  name           = "${var.table_name}"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "ID"

  attribute {
    name = "ID"
    type = "S"
  }

  tags {
    Environment = "AllDayDevOps"
  }
}
