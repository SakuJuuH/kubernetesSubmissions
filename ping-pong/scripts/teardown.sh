#!/bin/zsh

set -e

kubectl delete -f ./manifests/ -f ../log-output/manifests/ || true