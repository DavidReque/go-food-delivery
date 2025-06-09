#!/bin/bash

set -e

readonly service="$1"
readonly protoPath="api/protobuf/$service"

# Verificar si el directorio proto existe
if [ ! -d "$protoPath" ]; then
    echo "Skipping $service: Directory $protoPath does not exist"
    exit 0
fi

readonly outPath="./internal/services/$service/internal/shared/grpc/genproto"

# Crear el directorio de salida si no existe
mkdir -p "$outPath"

protoc \
  --proto_path="$protoPath" \
  --go_out="$outPath" \
  --go-grpc_out="$outPath" \
  --go-grpc_opt=require_unimplemented_servers=false \
  "$protoPath"/*.proto 