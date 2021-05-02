package data

import "database/sql"

type Models struct {
	Movies Movies
}

func New(db *sql.DB) Models  {
	return Models{Movies: MovieModel{DB: db}}
}