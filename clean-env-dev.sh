#!/bin/bash

echo "Removing Helm installation..."
helm uninstall kairos-ledger

echo "Cleaning up orphan images in the Kind cluster..."
docker exec -it kairos-vault-control-plane crictl rmi --prune

echo "Development environment cleaned!"