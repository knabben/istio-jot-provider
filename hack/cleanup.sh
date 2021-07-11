#!/usr/bin/env zsh

kubectl delete namespace istio-system
kubectl label namespace default istio-injection-
