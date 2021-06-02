resource "aws_lambda_function" "function" {
  function_name = "goidc-${var.name}"
  role          = aws_iam_role.execution.arn
  s3_bucket     = var.s3_bucket
  s3_key        = "sampleapp.zip"
  handler       = "main"
  runtime       = "go1.x"
  timeout       = 15
  publish       = true

  lifecycle {
    ignore_changes = [s3_bucket, s3_key]
  }

  depends_on = [aws_iam_role.execution]

  tags {
    goidc = yes
  }
}
