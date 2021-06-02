#!/bin/bash

bash -c ./scripts/ensure_aws.sh
if [[ $? -ne 0 ]]; then
    echo "Failed to setup AWS!"
    exit 1
fi

set -e

echo "Stage: $STAGE"
echo "Description: $DESCRIPTION"

echo "Deploying..."

aws apigateway create-deployment \
    --rest-api-id "$REST_API_ID" \
    --stage-name "$STAGE" \
    --description "$DESCRIPTION"

echo "Done!"