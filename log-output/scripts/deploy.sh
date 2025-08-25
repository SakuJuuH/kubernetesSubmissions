#!/bin/zsh

set -e

kubectl apply -f ./manifests/ -f ../ping-pong/manifests/