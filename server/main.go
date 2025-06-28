package main

import (
	"server/api"
)

func main() {
	router := api.NewRouter("/api")

	router.Get("/hi", func(reader *api.RequestReader) api.ResponseWriter {
		return api.SendOk("Hello world")
	})

	group := router.RouteGroup("/group")

	group.Get("/json", func(reader *api.RequestReader) api.ResponseWriter {
		return api.SendJson(true)
	})

	group.Get("/error", func(reader *api.RequestReader) api.ResponseWriter {
		return api.SendForbidden("You aren't allowed to do this")
	})

	if err := router.Listen(8080); err != nil {
		panic(err)
	}
}
