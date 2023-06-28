#! /bin/bash

# Resolve the version by the latest tag in the Git history, or fallback to "0.1.0".
VERSION=$(git describe --tags --abbrev=0 || echo "0.1.0")

# Build the Docker image with the resolved version.
docker build --build-arg "version=$VERSION" -t "cchantep/wilf:$VERSION" .
