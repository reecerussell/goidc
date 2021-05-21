#!/bin/bash

set -e

ensure_aws() {
    set +e
    aws --version &2>1 &1>/dev/null

    set -e

    if [[ $? -gt 0 ]]; then
        echo "Downloading AWS client..."
        curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
        unzip awscliv2.zip
        echo "Installing AWS client..."
        ./aws/install
    fi

    echo "Configuring AWS..."
    aws configure set default.region "$AWS_REGION"
    aws configure set aws_access_key_id "$AWS_ACCESS_KEY"
    aws configure set aws_secret_access_key "$AWS_SECRET_KEY"
}

echo "Publishing $FILE..."
echo "S3 Bucket: $BUCKET_NAME"
echo "S3 Key: $S3_KEY"

ensure_aws

echo "Copying to S3..."
aws s3 cp "$FILE" s3://$BUCKET_NAME/$S3_KEY

echo "Successfully published!"
