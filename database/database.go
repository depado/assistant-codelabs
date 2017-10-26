package database

import (
	"time"

	"github.com/Depado/assistant-codelabs/conf"
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
)

// DB is the main database. Put in separate package for use in external ones.
var DB *storm.DB

// Init initializes the database (creating it if necessary)
func Init() error {
	var err error
	DB, err = storm.Open(conf.C.DB, storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	return err
}
