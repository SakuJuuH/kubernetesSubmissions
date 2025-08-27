#!/bin/zsh

set -e

kubectl delete -k ./manifests || true