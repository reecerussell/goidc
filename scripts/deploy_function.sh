#!/bin/bash

bash -c ./scripts/ensure_aws.sh
if [[ $? -ne 0 ]]; then
    echo "Failed to setup AWS!"
    exit 1
fi

set -e

echo "Function: $NAME"
echo "Alias: $STAGE"
echo "Version: $VERSION"

echo "Updating alias..."

aws lambda update-alias \
    --function-name "$NAME" \
    --name "$STAGE" \
    --function-version "$VERSION"

echo "Done!"