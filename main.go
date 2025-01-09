package main

import (
	"database/sql"
	"log"

	"github.com/emmaahmads/summafy/api"
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel) // NEW
}

func main() {
	config, err := util.LoadConfigApp(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	util.MyGinLogger("ProductionMode...", config.ProductionMode)

	server := api.NewServer(*store, "superumi-summafy-123", config.SecretKey)
	err = server.Start( /* config.ServerAddress +  */ ":" + config.Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
