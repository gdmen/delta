package main

import (
	"database/sql"
	api "github.com/gdmen/delta/src/api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	os.Remove("./real.db")
	db, err := sql.Open("sqlite3", "./real.db")
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
