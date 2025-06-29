package main

import (
	"server/api"
)

func main() {
	router := api.NewRouter("/api")

	if err := router.Listen(8080); err != nil {
		panic(err)
	}
}
