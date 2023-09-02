package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nphattai/go-simple-bank/api"
	db "github.com/nphattai/go-simple-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address  = "0.0.0.0:80"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewSever(store)
	server.Start(address)
}
