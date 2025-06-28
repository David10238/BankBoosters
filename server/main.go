package main

import (
	"server/api"
)

func main() {
	router := api.NewRouter("/api")

	router.Get("/hi", func(reader *api.RequestReader) api.ResponseWriter {
		return api.SendOk("Hello world")
	})

	type Adder struct {
		A int
		B int
	}

	router.Post("/add", func(reader *api.RequestReader) api.ResponseWriter {
		adder := Adder{}

		if err := reader.BindJsonBody(&adder); err != nil {
			return err
		}

		d := 0
		if err := reader.BindJsonHeader("d", &d); err != nil {
			return err
		}

		str := ""
		if err := reader.BindStringHeader("hi", &str); err != nil {
			return err
		}

		if str == "Hi" {
			return api.SendJson("Hello world")
		}

		return api.SendJson(adder.A + adder.B + d)
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
