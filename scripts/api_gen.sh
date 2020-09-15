#!/bin/bash
if [ "$1" == "" ] || [ "$2" == "" ]; then
	echo "Proto file directory missing. Usage: ./scripts/api_gen.sh [subpath] [versionpath]"
	exit 1
fi

SUBPATH=$1
VERSION_PATH=$2

PROTOC_CMD="protoc  --proto_path=$GOPATH/src/github.com/huynguyen-quoc/go --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --go_out=plugins=grpc:$GOPATH/src $SERVICE_DIR/${SUBPATH}/pb/${VERSION_PATH}/*.proto"

$PROTOC_CMD

if [ "$?" -ne "0" ]; then
    echo -e "\033[1;31mCommand failed:\033[0m $PROTOC_CMD"
    exit 1
fi

echo "Compiled $PROTO_DIR"


echo "Finished running $0 script."
