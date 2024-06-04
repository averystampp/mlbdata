package main

import (
	"log"

	"github.com/averystampp/mlb"
	"github.com/averystampp/sesame"
	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("players"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("csrf"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("sessions"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("teams"))
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	db.Close()

	rtr := sesame.NewRouter()
	mlb.StartMLBService(rtr)
	rtr.StartServer(":5000")
}
