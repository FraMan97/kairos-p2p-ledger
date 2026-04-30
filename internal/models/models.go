package models

type AnchorRequest struct {
	Hash string `json:"hash"`
}

type LedgerRecord struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Ots       []byte `json:"ots"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}
