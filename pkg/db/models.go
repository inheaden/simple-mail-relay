package db

import "time"

// IP stores rate limiting information about an ip
type IP struct {
	IP           string
	LastCallTime time.Time
}

// Nonce stores information about which IP a nonce belongs to
type Nonce struct {
	IP       string
	Nonce    string
	SendTime time.Time
}
