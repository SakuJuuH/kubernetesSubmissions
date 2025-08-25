#!/bin/zsh

set -e

kubectl delete -f ./manifests/ -f ../manifests/ || true