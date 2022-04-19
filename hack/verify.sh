#!/bin/bash

echo "===( gofmt -d . )==="
fmt=$(gofmt -l .)
if [ "$fmt" != "" ]; then
    echo $fmt
    exit 1
fi
echo "===(     OK     )==="
echo ""

echo "===( go vet ./... )==="
go vet ./...
echo "===(      OK      )==="
echo ""

echo "===( go test ./... )==="
go test ./...
echo "===(      OK       )==="
echo ""

echo "===( protobuf verify )==="
echo "not implemented..."
