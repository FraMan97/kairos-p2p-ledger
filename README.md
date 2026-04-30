# Kairos P2P Ledger

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Helm](https://img.shields.io/badge/Helm-0F162D?style=for-the-badge&logo=helm&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![BadgerDB](https://img.shields.io/badge/BadgerDB-FFD700?style=for-the-badge&logo=go&logoColor=black)
![Bitcoin](https://img.shields.io/badge/Bitcoin-000?style=for-the-badge&logo=bitcoin&logoColor=white)

The **Kairos P2P Ledger** is the official tamper-evident notary service for the **kairos-p2p-engine** ecosystem. It provides a localized blockchain structure backed by BadgerDB and automatically anchors file proofs to the Bitcoin blockchain via OpenTimestamps, ensuring absolute, mathematically verifiable Proof of Existence.

> [!WARNING]
> This repository is for portfolio and demonstration purposes only. The source code is copyrighted and no license is granted for its use, modification, or distribution.
> This project is a Proof of Concept (POC) focused on backend infrastructure and is not intended for production environments without further security audits.

---

## Overview

It acts as an immutable audit log and Zero-Knowledge verification layer. By receiving only non-reversible hashes from the Gateway, it maintains absolute user privacy. The Ledger chains these hashes locally and periodically upgrades them into indisputable `.ots` cryptographic receipts.

### Key Features
* **Tamper-Evident Chain**: Securely links records using a localized blockchain structure (`PrevHash`) built on top of high-performance BadgerDB.
* **Free Bitcoin Anchoring**: Integrates natively with OpenTimestamps to batch and anchor proofs to the Bitcoin network without incurring gas fees.
* **Zero-Knowledge Architecture**: The service operates entirely blind. It receives and stores only SHA-256 hashes, ensuring zero exposure of user data or file contents.
* **Asynchronous Upgrades**: A background cron worker autonomously polls the OpenTimestamps calendar servers to upgrade pending receipts into fully offline-verifiable proofs.

## API Endpoints

The ledger exposes the following REST endpoints:

### Anchoring & Verification
* `POST /anchor`: Receives a SHA-256 hash, chains it to the local BadgerDB ledger, and submits it to the OpenTimestamps calendar server. Returns a `200 OK` when the hash is successfully registered as `PENDING`.
* `GET /receipt?hash=<hash>`: Retrieves the cryptographic proof for a given hash. 
  * Returns `202 Accepted` if the transaction is still pending on the blockchain.
  * Returns `200 OK` and initiates the download of the `.ots` binary file if the anchor is complete.

## Configuration

The ledger is configured via environment variables on the `values.yaml` helm file:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Local port the ledger listens on | `8080` |
| `DB_PATH` | Path where BadgerDB persists the local chain | `/data` |
| `OTS_DIGEST_URL` | Endpoint to submit new hashes | `https://a.pool.opentimestamps.org/digest` |
| `OTS_UPGRADE_URL` | Endpoint to upgrade pending receipts | `https://a.pool.opentimestamps.org/upgrade` |

## Getting Started

### Prerequisites
* **Go**: 1.24+
* **Kairos P2P Gateway**: A running instance of the gateway to route the hashing requests.

### Docker Deploy
To build the image and deploy the StatefulSet to your local Kubernetes cluster:
```bash
    ./deploy-dev.sh
```