package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/greenlight-api/internal/data"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint
// Now we just retrieve the interpolated "id" parameter from the current url
// and include it in a placeholder response
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Commando",
		Runtime:   120,
		Genre:     []string{"acappela", "war", "action"},
		Version:   12,
	}

	// otherwise, we interpolate the movie ID in a placeholder response
	// fmt.Fprintf(w, "show the details of movie %d\n", id)
	app.writeJson(data, w, r)
}
