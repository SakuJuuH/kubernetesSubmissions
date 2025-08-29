#!/bin/zsh

set -e

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <environment>"
    exit 1
fi

ENV=$1

kubectl apply -k ./kubernetes/overlays/"$ENV"