#!/bin/bash

set -e

echo "Building..."
echo "Working Directory: $WORKING_DIRECTORY"

cd "$WORKING_DIRECTORY"

export NODE_ENV=production
export CI=true

echo "Installing NPM modules..."
npm i
npm i --save-dev --unsafe-perm node-sass

echo "Building..."
npm run build

# Remove the prefixing slash on '/authorize', so that the paths
# map correctly, once in the API.
sed -i 's/\/authorize/authorize/g' build/index.html