package main

import (
	"hackaton-jam-back/views"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/cors"

	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func main() {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("HackatonJam API", "1.0.0"))

	// db := ConnectDB()
	// defer db.Close()

	handler := cors.AllowAll().Handler(router)

	views.Route(api)

	http.ListenAndServe("127.0.0.1:8888", handler)
}
