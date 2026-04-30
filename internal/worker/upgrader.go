package worker

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FraMan97/kairos-p2p-ledger/internal/config"
	"github.com/FraMan97/kairos-p2p-ledger/internal/database"
	"github.com/FraMan97/kairos-p2p-ledger/internal/models"
	"github.com/dgraph-io/badger/v4"
)

func StartUpgrader() {
	log.Printf("[Worker] Starting OTS upgrader cronjob. Interval: 1 hour")
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			processPendingRecords()
		}
	}()
}

func processPendingRecords() {
	log.Printf("[Worker] Running upgrade cycle for pending timestamps")

	var pendingRecords []models.LedgerRecord

	err := database.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var rec models.LedgerRecord
				if err := json.Unmarshal(v, &rec); err == nil {
					if rec.Status == "PENDING" {
						pendingRecords = append(pendingRecords, rec)
					}
				}
				return nil
			})
			if err != nil {
				log.Printf("[Worker] Error parsing item: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("[Worker] Error scanning database: %v", err)
		return
	}

	log.Printf("[Worker] Found %d pending records to process", len(pendingRecords))

	for _, rec := range pendingRecords {
		otsResp, err := http.Post(config.OtsUpgradeURL, "application/octet-stream", bytes.NewReader(rec.Ots))
		if err != nil {
			log.Printf("[Worker] Failed to contact upgrade pool for hash %s: %v", rec.Hash, err)
			continue
		}

		if otsResp.StatusCode == http.StatusOK {
			newOts, err := io.ReadAll(otsResp.Body)
			otsResp.Body.Close()

			if err == nil && len(newOts) > 0 && len(newOts) != len(rec.Ots) {
				rec.Ots = newOts
				rec.Status = "ANCHORED"

				recordBytes, _ := json.Marshal(rec)
				database.Put(rec.Hash, recordBytes)
				log.Printf("[Worker] Successfully upgraded and anchored hash %s", rec.Hash)
			} else {
				log.Printf("[Worker] Hash %s is still pending consolidation", rec.Hash)
			}
		} else {
			otsResp.Body.Close()
		}
	}
}
