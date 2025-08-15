#!/bin/bash

# Kiểm tra tham số đầu vào
if [ "$1" = "swag" ]; then
  echo "Usage: $0 swag [extra_flags]"
  cd ../../$2-service && swag init --parseDependency --output ../../gateway-service/docs/$2

  exit 1
fi

if [ "$1" = "wire" ]; then
  echo "Usage: $1 swag [extra_flags]"
  cd ../../$2-service
  echo "$(pwd)"
  go run ../shared/code/script/gen_wire.go $2

  exit 1
fi


if [ "$1" = "social" ]; then
  echo "Usage: $1 swag [extra_flags]"
  cd social-service 
  swag init --parseDependency --output ../../gateway-service/docs/social
  go run script/gen_wire.go social
  exit 1
fi

