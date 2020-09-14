#!/bin/bash

if [ "$1" == "" ]; then
	echo "Proto file directory missing. Usage: ./scripts/stream_gen.sh [subpath]"
	exit 1
fi

SUBPATH=$1

# Formatting stream_info.proto to correct package
PROTOC_CMD="protoc --proto_path=$STREAM_PROTO_DIR --proto_path=$STREAM_PROTO_DIR/streams/common --proto_path=$GOPATH/src/ --go_out=plugins=grpc:$GOPATH/src $STREAM_PROTO_DIR/streams/${SUBPATH}/*.proto"

$PROTOC_CMD

if [ "$?" -ne "0" ]; then
    echo -e "\033[1;31mCommand failed:\033[0m $PROTOC_CMD"
    echo -e "\033[1;31mCheck the error messages to see if there's a proto syntax failure."
    echo -e "If you get some Import was not found error, you likely missed some files in /usr/local/include/"
    echo -e "follow the readme from the downloaded protoc release zip file https://github.com/protocolbuffers/protobuf/releases/tag/v3.5.1"
    echo -e "Otherwise please check the --proto_path and see if there's an obvious fix."
    echo -e "If you're still stuck, please contact the authors named in the script header"
    echo -e "or check the commit logs for recent contributors.\033[0m "
    exit 1
fi

echo "Compiled $PROTO_DIR"


echo "Finished running $0 script."
