package main

import (
	"database/sql"
	"log"

	"github.com/emmaahmads/summafy/api"
	db "github.com/emmaahmads/summafy/db/sqlc"
	_ "github.com/emmaahmads/summafy/docs"
	"github.com/emmaahmads/summafy/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel) // NEW
}

//	@title			Swagger Example API emma
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
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
