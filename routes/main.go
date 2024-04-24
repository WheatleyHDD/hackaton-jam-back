package routes

import (
	"database/sql"

	"hackaton-jam-back/routes/auth"
	"hackaton-jam-back/routes/example"

	"github.com/danielgtaylor/huma/v2"

	_ "github.com/lib/pq"
)

func Route(api huma.API, db *sql.DB) {
	example.Route(api, db)
	auth.Route(api, db)
}
