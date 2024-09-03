package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/emmaahmads/summafy/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	TestDB, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(TestDB)
	os.Exit(m.Run())
}
