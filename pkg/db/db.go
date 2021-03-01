package db

import (
	"time"

	"github.com/apex/log"
	"github.com/hashicorp/go-memdb"
)

type Database struct {
	db *memdb.MemDB
}

func NewDB() *Database {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"ip": {
				Name: "ip",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "IP"},
					},
				},
			},
			"nonce": {
				Name: "nonce",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Nonce"},
					},
					"ip": {
						Name:    "ip",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "IP"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return &Database{db: db}
}

// GetIP returns an ip or nil
func (d *Database) GetIP(ip string) (*IP, error) {
	txn := d.db.Txn(false)
	defer txn.Commit()

	result, err := txn.First("ip", "id", ip)
	if result == nil || err != nil {
		return nil, err
	}
	return result.(*IP), nil
}

// SaveIP saves or updates an ip
func (d *Database) SaveIP(ip string) error {
	txn := d.db.Txn(true)
	defer txn.Commit()

	return txn.Insert("ip", &IP{IP: ip, LastCallTime: time.Now()})
}

// GetNonce returns a nonce or nil
func (d *Database) GetNonce(nonce string) (*Nonce, error) {
	txn := d.db.Txn(false)
	defer txn.Commit()

	result, err := txn.First("nonce", "id", nonce)
	if result == nil || err != nil {
		return nil, err
	}
	return result.(*Nonce), err
}

// SaveNonce saves or updates a nonce
func (d *Database) SaveNonce(nonce string, ip string) error {
	txn := d.db.Txn(true)
	defer txn.Commit()

	return txn.Insert("nonce", &Nonce{Nonce: nonce, IP: ip, SendTime: time.Now()})
}

// DeleteNonceByNonce deletes a nonce using its ip
func (d *Database) DeleteNonceByNonce(nonce string) error {
	txn := d.db.Txn(true)
	defer txn.Commit()

	nonceModel, err := d.GetNonce(nonce)
	if nonceModel == nil || err != nil {
		return err
	}
	return txn.Delete("nonce", nonceModel)
}

// CleanIPs removes all ip information that is older than maxAgeSeconds.
func (d *Database) CleanIPs(maxAgeSeconds int) error {
	txn := d.db.Txn(true)
	defer txn.Commit()

	it, err := txn.Get("ip", "id")
	if err != nil {
		return err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		ip := obj.(*IP)
		if ip.LastCallTime.Before(time.Now().Add(-time.Duration(maxAgeSeconds) * time.Second)) {
			err = txn.Delete("ip", ip)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanNonces removes all nonce information that is older than maxAgeSeconds.
func (d *Database) CleanNonces(maxAgeSeconds int) error {
	txn := d.db.Txn(true)
	defer txn.Commit()

	it, err := txn.Get("nonce", "id")
	if err != nil {
		return err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		nonce := obj.(*Nonce)
		if nonce.SendTime.Before(time.Now().Add(-time.Duration(maxAgeSeconds) * time.Second)) {
			log.Infof("%v", nonce)
			err = txn.Delete("nonce", nonce)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
