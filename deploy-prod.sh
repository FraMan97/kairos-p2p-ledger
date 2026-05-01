#!/bin/bash

helm upgrade --install kairos-ledger ./helm \
  -f ./helm/values-dev.yaml \
  -f ./helm/values-prod.yaml \
  --namespace production --create-namespace

kubectl rollout restart statefulset/kairos-ledger -n production

echo "Deploy Ledger completed!"