package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nphattai/go-simple-bank/util"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE)

	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
