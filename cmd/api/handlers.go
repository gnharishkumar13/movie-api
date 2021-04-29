package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//js := `{"status": "available", "environment": %q, "version": %q}`
	//js = fmt.Sprintf(js, app.config.env, version)
	//w.Header().Set("Content-Type", "application/json")

	data := map[string]string{
		"status":"available",
		"environment": app.config.env,
		"version": version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}


func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Otherwise, interpolate the movie ID in a placeholder response.
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}