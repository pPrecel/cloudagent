#!/bin/bash

export GIT_REF="local"
if [ "$GITHUB_REF_NAME" != "" ]; then
    GIT_REF="$GITHUB_REF_NAME"
fi

echo "building cloudagent with version: $GIT_REF"
go build -ldflags="-X 'main.Version=${GIT_REF}'" -o .out/cloudagent main.go
