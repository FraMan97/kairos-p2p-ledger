package database

import (
	"encoding/json"
	"log"
	"os"

	"github.com/FraMan97/kairos-p2p-ledger/internal/models"
	"github.com/dgraph-io/badger/v4"
)

var DB *badger.DB

const HeadKey = "_CHAIN_HEAD_"

func InitDB(dbPath string) {
	log.Printf("[Database] Initializing BadgerDB at %s", dbPath)

	err := os.MkdirAll(dbPath, 0700)
	if err != nil {
		log.Fatalf("[Database] Failed to create directory: %v", err)
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("[Database] Critical failure opening BadgerDB: %v", err)
	}

	DB = db
	log.Printf("[Database] BadgerDB successfully opened")
}

func AppendRecord(record *models.LedgerRecord) error {
	return DB.Update(func(txn *badger.Txn) error {
		var prevHash string
		item, err := txn.Get([]byte(HeadKey))
		if err == nil {
			val, _ := item.ValueCopy(nil)
			prevHash = string(val)
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		record.PrevHash = prevHash

		recordBytes, err := json.Marshal(record)
		if err != nil {
			return err
		}

		if err := txn.Set([]byte(record.Hash), recordBytes); err != nil {
			return err
		}

		return txn.Set([]byte(HeadKey), []byte(record.Hash))
	})
}

func Put(key string, val []byte) error {
	return DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), val)
	})
}

func Get(key string) ([]byte, error) {
	var valCopy []byte
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(nil)
		return err
	})
	return valCopy, err
}

func CloseDB() {
	if DB != nil {
		log.Printf("[Database] Closing BadgerDB safely")
		DB.Close()
	}
}
