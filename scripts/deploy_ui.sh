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

aws apigateway update-stage \
    --rest-api-id "$REST_API_ID" \
    --stage-name "$ENV" \
    --patch-operations op="replace",path=/variables/"$VAR",value=$VERSION