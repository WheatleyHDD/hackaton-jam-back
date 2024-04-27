package main

import (
	"fmt"
	"hackaton-jam-back/routes"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/cors"

	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port       int    `help:"Port to listen on" short:"p" default:"8888"`
	DbPassword string `help:"Database password" default:"default"`
	Ip         string `help:"Ip address to listen on" short:"i" default:"127.0.0.1"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := http.NewServeMux()
		api := humago.New(router, huma.DefaultConfig("HackatonJam API", "1.0.0"))

		db := ConnectDB(options.DbPassword)

		handler := cors.AllowAll().Handler(router)

		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, fmt.Sprintf("http://%s/docs", r.Host), http.StatusFound)
		})
		routes.Route(api, db)

		hooks.OnStart(func() {
			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", options.Ip, options.Port), handler); err != nil {
				log.Fatalf("HTTP server error: %v", err)
			}
			defer db.Close()
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
