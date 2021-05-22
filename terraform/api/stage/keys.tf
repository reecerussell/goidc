resource "aws_kms_key" "jwt" {
    description = "A signing key for the JWT handlers (${var.name})."
    key_usage = "SIGN_VERIFY"
    customer_master_key_spec = "RSA_2048"

    tags = {
        stage = var.name
    }
}