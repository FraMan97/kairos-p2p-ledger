package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FraMan97/kairos-p2p-ledger/internal/config"
	"github.com/FraMan97/kairos-p2p-ledger/internal/database"
	"github.com/FraMan97/kairos-p2p-ledger/internal/models"
)

func AnchorHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[API: AnchorHash] Invalid method %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AnchorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[API: AnchorHash] Failed to decode request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hashBytes, err := hex.DecodeString(req.Hash)
	if err != nil || len(hashBytes) != 32 {
		log.Printf("[API: AnchorHash] Invalid SHA256 hash provided: %s", req.Hash)
		http.Error(w, "Invalid SHA256 hash", http.StatusBadRequest)
		return
	}

	otsResp, err := http.Post(config.OtsDigestURL, "application/octet-stream", bytes.NewReader(hashBytes))
	if err != nil {
		log.Printf("[API: AnchorHash] Failed to contact OpenTimestamps pool: %v", err)
		http.Error(w, "Failed to contact timestamp pool", http.StatusBadGateway)
		return
	}
	defer otsResp.Body.Close()

	if otsResp.StatusCode != http.StatusOK {
		log.Printf("[API: AnchorHash] OpenTimestamps returned status %d", otsResp.StatusCode)
		http.Error(w, "Timestamp pool error", http.StatusBadGateway)
		return
	}

	otsData, err := io.ReadAll(otsResp.Body)
	if err != nil {
		log.Printf("[API: AnchorHash] Failed to read OpenTimestamps response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	record := models.LedgerRecord{
		Hash:      req.Hash,
		Ots:       otsData,
		Status:    "PENDING",
		Timestamp: time.Now().Unix(),
	}

	if err := database.AppendRecord(&record); err != nil {
		log.Printf("[API: AnchorHash] Database write failure for hash %s: %v", req.Hash, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("[API: AnchorHash] Successfully anchored hash %s. Status: PENDING", req.Hash)
	w.WriteHeader(http.StatusOK)
}

func GetReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("[API: GetReceipt] Invalid method %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		log.Printf("[API: GetReceipt] Missing hash parameter")
		http.Error(w, "Missing hash parameter", http.StatusBadRequest)
		return
	}

	data, err := database.Get(hash)
	if err != nil {
		log.Printf("[API: GetReceipt] Hash %s not found in database", hash)
		http.Error(w, "Hash not found", http.StatusNotFound)
		return
	}

	var record models.LedgerRecord
	if err := json.Unmarshal(data, &record); err != nil {
		log.Printf("[API: GetReceipt] Failed to unmarshal record for hash %s: %v", hash, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if record.Status != "ANCHORED" {
		log.Printf("[API: GetReceipt] Receipt for hash %s requested but status is %s", hash, record.Status)
		http.Error(w, "Receipt is still pending on the blockchain", http.StatusAccepted)
		return
	}

	log.Printf("[API: GetReceipt] Serving completed receipt for hash %s", hash)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.ots"`, hash))
	w.Write(record.Ots)
}
