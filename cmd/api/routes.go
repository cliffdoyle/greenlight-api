package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//Encapsulating our routing rules like this keeps our main() function clean and
//ensures all our routes are defined in one single place

func (app *application) routes() http.Handler {
	//Initialize a new httprouter router instance
	router := httprouter.New()

	//Register the relevant methods, URL patterns and handler functions
	//for our endpoints using the HandlerFunc() method.
	//http.MethodGet and http.MethodPost are constants which equate to the
	//strings "GET" and "POST" respectively

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	// Add the route for the PUT /v1/movies/:id endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateMovieHandler)

	//Return the httprouter instance
	return router
}
