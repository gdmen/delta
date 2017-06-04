package main

import (
	api "github.com/gdmen/delta/src/api"
)

func main() {
	api_router := api.GetRouterV1()
	api_router.Run(":8080")
}
