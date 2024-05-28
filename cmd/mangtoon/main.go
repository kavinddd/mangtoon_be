package main

import (
	"database/sql"
	"github.com/kavinddd/mangtoon_be/internal/db"
	"github.com/kavinddd/mangtoon_be/internal/rest"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	// load config fron .env
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	// open db-connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	// inject the connected db to run the server
	store := db.NewStore(conn)
	s, err := rest.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	// run the server
	if err := s.Run(config.ServerAddress); err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
