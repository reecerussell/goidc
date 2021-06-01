#!/bin/bash

bash -c ./scripts/ensure_aws.sh
if [[ $? -ne 0 ]]; then
    echo "Failed to setup AWS!"
    exit 1
fi

set -e

echo "Env: $ENV"
echo "Version: $VERSION"
echo "Variable: $VAR"
echo "Name: $NAME"

echo "Updating stage..."
aws apigateway update-stage \
    --rest-api-id "$REST_API_ID" \
    --stage-name "$ENV" \
    --patch-operations op="replace",path=/variables/"$VAR",value=$VERSION

echo "Deploying..."
aws apigateway create-deployment \
    --rest-api-id "$REST_API_ID" \
    --stage-name "$ENV" \
    --description "Deployed $NAME ($VERSION)"

echo "Done!"