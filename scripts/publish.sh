#!/bin/bash

bash -c ./scripts/ensure_aws.sh
if [[ $? -ne 0 ]]; then
    echo "Failed to setup AWS!"
    exit 1
fi

set -e

echo "Publishing $NAME ($FILE)..."
echo "S3 Bucket: $S3_BUCKET"
echo "S3 Key: $S3_KEY"

echo "Copying to S3..."
aws s3 cp "$FILE" s3://$S3_BUCKET/$S3_KEY

echo "Updating function..."
aws lambda update-function-code \
    --function-name "$NAME" \
    --s3-bucket "$S3_BUCKET" \
    --s3-key "$S3_KEY"

echo "Publish new version..."
aws lambda publish-version --function-name "$NAME" > /tmp/version

version=$(cat /tmp/version | jq '.Version' | sed 's/\"//g')
echo "::set-output name=version::$version"

echo "Successfully published!"
