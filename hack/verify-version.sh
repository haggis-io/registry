#!/bin/bash

if [[ "$LATEST_TAG" == "" ]]; then
    exit 1
elif [[ "$LATEST_TAG" = "$NEW_VERSION" ]]; then
    echo 'Version is already exists, please change!'
    exit 1
fi