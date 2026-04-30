#!/bin/bash

# 0. Kill any hanging processes occupying port 8087
echo "Cleaning up previous port-forwards..."
fuser -k 8087/tcp || true

# 1. Build the images
./build-images.sh

# 2. Load the images into the local Kind cluster
echo "Loading images into Kind..."
kind load docker-image kairos-ledger:local --name kairos-vault

# 3. Install the P2P infrastructure with Helm
echo "Installing Helm Chart..."
helm upgrade --install kairos-ledger -f ./helm/values-dev.yaml ./helm \
  --set ledger.image=kairos-ledger:local \
  --set ledger.pullPolicy=Never 

kubectl rollout restart statefulset/kairos-ledger

# 4. Wait for Kubernetes to finish the job
echo "Waiting for all pods to be ready (this may take a few seconds)..."
kubectl rollout status statefulset/kairos-ledger --timeout=120s

# 5. Configure Port Forwards (in background)
echo "Configuring Ledger API Port-Forward on localhost:8087..."
kubectl port-forward svc/kairos-ledger 8087:80 > /dev/null 2>&1 &

echo "Deploy completed! The Ledger is running."
echo "- Ledger API is on http://localhost:8087"