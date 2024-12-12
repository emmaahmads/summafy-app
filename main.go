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
	awsConf, err := util.LoadConfigAws(".", config.ProductionMode)
	if err != nil {
		log.Fatalf("cannot load AWS config: %v", err)
	}

	server := api.NewServer(*store, api.NewAwsConfig(
		awsConf.S3Bucket,
		awsConf.Region,
		awsConf.Creds1,
		awsConf.Creds2,
		awsConf.Creds3,
		awsConf.ApiKey,
	))
	err = server.Start(/* config.ServerAddress +  */":" + config.Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
