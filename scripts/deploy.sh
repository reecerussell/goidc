#!/bin/bash

set -e

echo "Deploying $NAME..."
echo "Source: $SOURCE_ZIP"

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

    if [[ "$AWS_REGION" != "" ]]; then
        aws configure set default.region "$AWS_REGION"
    fi

    if [[ "$AWS_ACCESS_KEY" != "" ]]; then
        aws configure set aws_access_key_id "$AWS_ACCESS_KEY"
    fi
    
    if [[ "$AWS_SECRET_KEY" != "" ]]; then
        aws configure set aws_secret_access_key "$AWS_SECRET_KEY"
    fi
}

echo "Publishing $SOURCE_ZIP..."
echo "Region: $AWS_REGION"
echo "S3 Bucket: $S3_BUCKET"
echo "Source: $SOURCE_ZIP"

ensure_aws

echo "Updating function..."
aws lambda update-function-code \
    --function-name "$NAME" \
    --s3-bucket "$S3_BUCKET" \
    --s3-key "$SOURCE_ZIP"

echo "Publishing new version..."
version=$(aws lambda publish-version --function-name "$NAME" | jq '.Version' | sed 's/"//g')

echo "Updating '$ENV' version..."
aws lambda update-alias \
    --function-name "$NAME" \
    --name "$ENV" \
    --function-version "$version"

echo "Deploying API..."
aws apigateway create-deployment \
    --rest-api-id "$REST_API_ID" \
    --stage-name "$ENV" \
    --description "Deployed $NAME ($SOURCE_ZIP)"

echo "Successfully Deployed!"
