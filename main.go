package main

import (
	"database/sql"
	"fmt"
	"hackaton-jam-back/routes"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/cors"
	"github.com/tanimutomo/sqlfile"

	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int    `help:"Port to listen on" short:"p" default:"8888"`
	Ip   string `help:"Ip address to listen on" short:"i" default:"127.0.0.1"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := http.NewServeMux()
		api := humago.New(router, huma.DefaultConfig("HackatonJam API", "1.0.0"))

		db := ConnectDB()

		handler := cors.AllowAll().Handler(router)

		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, fmt.Sprintf("http://%s/docs", r.Host), http.StatusFound)
		})
		routes.Route(api, db)

		hooks.OnStart(func() {
			if os.Getenv("FIRST_RUN") == "1" {
				fmt.Println("Started migration...")
				Migration(db)
				defer db.Close()
				fmt.Println("Migration ended!")
			} else {
				if err := http.ListenAndServe(fmt.Sprintf("%s:%d", options.Ip, options.Port), handler); err != nil {
					log.Fatalf("HTTP server error: %v", err)
				}
				defer db.Close()
			}
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
func Migration(db *sql.DB) {
	_ = db.QueryRow("DROP SCHEMA public CASCADE;" +
		"CREATE SCHEMA public;")

	s := sqlfile.New()

	// Load input file and store queries written in the file
	err := s.File("sql/hjam.sql")
	if err != nil {
		log.Fatal("Невозможно получить файл")
	}
	_, err = s.Exec(db)
	if err != nil {
		log.Fatal("Невозможно мигрировать базу данных")
	}
}
