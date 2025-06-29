package main

import (
	"server/api"
)

func main() {
	router := api.NewRouter("/api")

	if err := router.ListenAndServeTLS(8080, "ca.crt", "ca.key"); err != nil {
		panic(err)
	}
}
