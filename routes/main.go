package routes

import (
	"database/sql"

	"hackaton-jam-back/routes/auth"
	_ "hackaton-jam-back/routes/events"
	"hackaton-jam-back/routes/example"
	"hackaton-jam-back/routes/profile"

	"github.com/danielgtaylor/huma/v2"

	_ "github.com/lib/pq"
)

func Route(api huma.API, db *sql.DB) {
	example.Route(api, db)
	auth.Route(api, db)
	profile.Route(api, db)
}
