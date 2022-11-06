#!/bin/bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.2

KABU_STATION_DOC="https://kabucom.github.io/kabusapi/reference/kabu_STATION_API.yaml"
KABU_STATION_PACKAGE="../entities/infrastructures/kabucom/autogen"

cd "$(dirname "$0")"
mkdir -p $KABU_STATION_PACKAGE
oapi-codegen -generate spec,types,client -package autogen $KABU_STATION_DOC > "$KABU_STATION_PACKAGE/client.gen.go"