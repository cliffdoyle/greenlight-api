package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/greenlight-api/internal/data"
	"github.com/greenlight-api/validator"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
		Version int32    `json:"version"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	//Make the movie variable contains a pointer to a movie struct
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	data.ValidateMovie(v, movie)
	if !v.Valid() {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	//Call in the Insert() method passing in a pointer to the validated movie struct
	//This will create a record in the database and update the
	//movie struct with the system-generated information
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//When sending a HTTP response, we want to include a Location header to let the
	//client know which URL they can find the newly-created resource at.
	//we make an empty http.Header map and then use the Set() method to add a new Location header,
	//interpolating the system-generated ID for our new movie in the URL
	headers := make(http.Header)

	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	//Write a JSON response with a 201 created status code, the movie data in the
	//response body, and the Location header

	app.writeJson(map[string]any{"movie": movie}, w, r)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint
// Now we just retrieve the interpolated "id" parameter from the current url
// and include it in a placeholder response
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//call the Get() method to fetch the data for a specific movie. We also need to use
	//the errors.Is() function to check if it returns a data.ErrRecordNotFound error
	//in which case we send a 404 not found response to the client

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	app.writeJson(map[string]any{"movie": movie}, w, r)
}
