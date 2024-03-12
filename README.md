# hclwrite-dynamodb

PoC hclwrite package

Terraform連載2024 技術ブログ記事に用いたサンプルコード。

## Example

```sh
$ go run . example/dynamodb_table.tf
// DO NOT EDIT, MADE BY hclwrite-dynamodb-generator

resource "aws_cloudwatch_metric_alarm" "dynamodb_throttledrequests_myproduct_read" {
  alarm_name          = "${aws_dynamodb_table.myproduct_read.name}-throttledrequests"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  datapoints_to_alarm = "1"
  evaluation_periods  = "1"
  metric_name         = "ThrottledRequests"
  namespace           = "AWS/DynamoDB"
  period              = "60"
  statistic           = "Maximum"
  threshold           = "1"
  alarm_actions       = aws_sns_topic.myproduct_alert.arn
  dimensions {
    TableName = aws_dynamodb_table.myproduct_read.name
  }
}

resource "aws_cloudwatch_metric_alarm" "dynamodb_throttledrequests_myproduct_content" {
  alarm_name          = "${aws_dynamodb_table.myproduct_read.name}-throttledrequests"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  datapoints_to_alarm = "1"
  evaluation_periods  = "1"
  metric_name         = "ThrottledRequests"
  namespace           = "AWS/DynamoDB"
  period              = "60"
  statistic           = "Maximum"
  threshold           = "1"
  alarm_actions       = aws_sns_topic.myproduct_alert.arn
  dimensions {
    TableName = aws_dynamodb_table.myproduct_read.name
  }
}
```
