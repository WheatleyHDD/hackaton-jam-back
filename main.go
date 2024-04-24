package main

import (
	"hackaton-jam-back/routes"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/cors"

	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func main() {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("HackatonJam API", "1.0.0"))

	db := ConnectDB()
	defer db.Close()

	handler := cors.AllowAll().Handler(router)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Location", "/docs")
	})
	routes.Route(api, db)

	http.ListenAndServe(":8888", handler)
}
