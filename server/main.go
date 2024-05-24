package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/averystampp/mlb"
	"github.com/averystampp/sesame"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	connStr := "postgresql://postgres:docker@localhost:5432/postgres?sslmode=disable"
	sesame.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rtr := sesame.NewRouter()
	players := flag.Bool("players", false, "")
	flag.Parse()
	if *players {
		err = mlb.GetAllPlayers()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Got all player data")
	} else {
		fmt.Println("Starting without player pull")
	}
	mlb.StartMLBService(rtr)

	rtr.StartServer(":5000")
}
