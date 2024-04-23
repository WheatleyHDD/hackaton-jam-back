package views

import (
	"database/sql"

	"github.com/danielgtaylor/huma/v2"

	_ "github.com/lib/pq"
)

type ApiRoutes struct {
	Api huma.API
	Db  *sql.DB
}

// db *sql.DB
func Route(api huma.API) {
	apiRoutes := ApiRoutes{Api: api}

	apiRoutes.Example()
}
