#!/bin/zsh

set -e

if [ $# -ne 1 ]; then
  echo "Usage: $0 <tag>"
  exit 1
fi

TAG=$1
IMAGE_NAME="sakuheinonen/$(basename $(pwd))"

docker build --platform=linux/amd64 -t ${IMAGE_NAME}:${TAG} .
docker push ${IMAGE_NAME}:${TAG}