#!/usr/bin/env zsh

ISTIO_VERSION="1.10.2"

# Download istio local
curl -L https://istio.io/downloadIstio | sh -
export PATH="$PATH:${PWD}/istio-${ISTIO_VERSION}/bin"

# Install profile setup
istioctl x precheck
istioctl install --set profile=demo -y

# Set namespace label for istio usage
kubectl label namespace default istio-injection=enabled