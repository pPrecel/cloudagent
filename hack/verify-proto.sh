#!/bin/bash

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

git diff --exit-code > /dev/null
if [ $? -eq 1 ]; then
    echo "git in dirty state"
    exit 1
fi

echo ">==( protobuf verify )==="

trap "git reset --hard" EXIT

make -C "${SCRIPTPATH}/.." protobuf

if [ "$(git diff)" != "" ]; then
    echo "protobuf is not up to date"
    exit 1
fi

echo "<==(       OK        )==="
