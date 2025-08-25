#!/bin/zsh

set -e

kubectl delete -f ./manifests/ -f ../ping-pong/manifests/ || true