package main

import (
	"database/sql"
	"log"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	_ = db.NewStore(conn)

}
