resource "aws_s3_bucket" "ui_bucket" {
  bucket = "goidc-authorize-ui"
  acl    = "private"
}
