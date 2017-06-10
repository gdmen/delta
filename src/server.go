package main

import (
	"database/sql"
	api "github.com/gdmen/delta/src/api"
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

	api := api.NewApi(db)
	api_router := api.GetRouter()
	api_router.Run(":8080")
}
