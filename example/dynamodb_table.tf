resource "aws_dynamodb_table" "myproduct_read" {
  name         = "${terraform.workspace}-myproduct-read"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "user_id"
  range_key    = "content_id"

  deletion_protection_enabled = true

  attribute {
    name = "user_id"
    type = "S"
  }
  attribute {
    name = "content_id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "myproduct_content" {
  name         = "${terraform.workspace}-myproduct-read"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "content_id"

  deletion_protection_enabled = true

  attribute {
    name = "content_id"
    type = "S"
  }
}