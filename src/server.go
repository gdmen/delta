package main

import (
	"database/sql"
	"fmt"
	"github.com/gdmen/delta/src/api"
	"github.com/gdmen/delta/src/common"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	c, err := common.ReadConfig("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", c.MySQLUser, c.MySQLPass, c.MySQLHost, c.MySQLPort, c.MySQLDatabase)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api, err := api.NewApi(db)
	if err != nil {
		log.Fatal(err)
	}
	api_router := api.GetRouter()
	api_router.Run(":8080")
}
