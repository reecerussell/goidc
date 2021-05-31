#!/bin/bash

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

if [[ "$AWS_REGION" != "" ]]; then
    aws configure set default.region "$AWS_REGION"
fi

if [[ "$AWS_ACCESS_KEY" != "" ]]; then
    aws configure set aws_access_key_id "$AWS_ACCESS_KEY"
fi

if [[ "$AWS_SECRET_KEY" != "" ]]; then
    aws configure set aws_secret_access_key "$AWS_SECRET_KEY"
fi

exit 0