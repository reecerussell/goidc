resource "aws_lambda_function" "function" {
  function_name = var.name
  role          = var.execution_role
  s3_bucket     = var.s3_bucket
  s3_key        = "sampleapp.zip"
  handler       = "main"
  runtime       = "go1.x"
  timeout       = 15
  publish       = true

  lifecycle {
    ignore_changes = [s3_bucket, s3_key]
  }
}
