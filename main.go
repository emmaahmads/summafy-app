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
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	aws_conf := api.NewAwsConfig(config.S3Bucket, config.Region, config.Creds1, config.Creds2, config.Creds3)
	server := api.NewServer(*store, aws_conf, config.ApiKey)
	err = server.Start(config.ServerAddress + ":" + config.Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
