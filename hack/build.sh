#!/bin/bash

GIT_REF="$(git describe --exact-match --tags $(git log -n1 --pretty='%h') 2>/dev/null)"
STATUS="$?"

if [ "$STATUS" != "0" ]; then
    export GIT_REF="$(git branch --show-current)"
fi

go build -ldflags="-X 'main.Version=${GIT_REF}'" -o .out/cloudagent main.go
