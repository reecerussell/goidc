resource "aws_iam_role" "execution" {
  name = "${var.name}-exec"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "logging" {
  name        = "${var.name}-logging"
  path        = "/"
  description = "IAM policy for logging for ${var.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "logs_attachment" {
  role       = aws_iam_role.execution.name
  policy_arn = aws_iam_policy.logging.arn

  depends_on = [aws_iam_role.execution]
}

resource "aws_iam_role_policy_attachment" "attachment" {
  count = length(var.iam_policies)

  role       = aws_iam_role.execution.name
  policy_arn = var.iam_policies[count.index]

  depends_on = [aws_iam_role.execution]
}