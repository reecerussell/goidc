#!/bin/bash

bash -c ./scripts/ensure_aws.sh
if [[ $? -ne 0 ]]; then
    echo "Failed to setup AWS!"
    exit 1
fi

set -e

echo "S3 Bucket: $S3_BUCKET"
echo "Version: $VERSION"

build_dir=$BUILD_DIR
if [[ $build_dir = "" ]]; then
    build_dir="build/."
fi

echo "Build Directory: $build_dir"

echo "Syncing Files..."
aws s3 sync "$build_dir" s3://$S3_BUCKET/$VERSION/

echo "Done!"
