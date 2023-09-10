package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nphattai/go-simple-bank/api"
	db "github.com/nphattai/go-simple-bank/db/sqlc"
	"github.com/nphattai/go-simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config :", err)
	}

	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)

	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewSever(store)

	if err != nil {
		log.Fatal("can not create server: ", err)
	}
	server.Start(config.SERVER_ADDRESS)
}
