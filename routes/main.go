package routes

import (
	"database/sql"

	"hackaton-jam-back/routes/auth"
	"hackaton-jam-back/routes/events"
	"hackaton-jam-back/routes/example"
	"hackaton-jam-back/routes/notifications"
	"hackaton-jam-back/routes/profile"
	"hackaton-jam-back/routes/teams"

	"github.com/danielgtaylor/huma/v2"

	_ "github.com/lib/pq"
)

func Route(api huma.API, db *sql.DB) {
	example.Route(api, db)
	auth.Route(api, db)
	profile.Route(api, db)
	events.Route(api, db)
	notifications.Route(api, db)
	teams.Route(api, db)
}
