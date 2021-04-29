package main

import (
	"fmt"
	"github.com/gnharishkumar13/movie-api/internal/data"
	"net/http"
	"time"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//js := `{"status": "available", "environment": %q, "version": %q}`
	//js = fmt.Sprintf(js, app.config.env, version)
	//w.Header().Set("Content-Type", "application/json")

	data := envelope{
		"status":"available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	m := data.Movie{
		ID:        id,
		CreatedAt: time.Time{},
		Title:     "test",
		Year:      2021,
		Runtime:   150,
		Genres:    []string{"test","test1"},
		Version:   20,

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": m}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}