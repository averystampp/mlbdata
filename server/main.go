package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/averystampp/mlb"
	"github.com/averystampp/sesame"
	bolt "go.etcd.io/bbolt"
)

func main() {
	// var err error
	// connStr := "postgresql://postgres:docker@localhost:5432/postgres?sslmode=disable"
	// sesame.DB, err = sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
	players := flag.Bool("players", false, "")
	flag.Parse()
	if *players {
		err := mlb.GetAllPlayers()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Got all player data")
	} else {
		fmt.Println("Starting without player pull")
	}
	mlb.StartMLBService(rtr)
	mlb.DraftService(rtr)
	mlb.UserService(rtr)
	rtr.StartServer(":5000")
}
