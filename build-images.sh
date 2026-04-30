#!/bin/bash

echo "Building Ledger image..."
docker build -t kairos-ledger:local -f Dockerfile.ledger .