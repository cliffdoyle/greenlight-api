package main

import (
	"fmt"
	"net/http"
)

// The logError() method is ageneric helper for logging an error message
// along with the current request method and URL as attribute in the log entry
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse() Method is a generic helper for sending JSON-formatted error messages
// to the client with a given status code.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	data := map[string]any{"error": message}

	//Write the response using the writeJSON() helper.If this happens to
	//return an error then we log it, and fall back to sending the client
	//an empty response with a 500 internal server error status code
	app.writeJson(data,w,r)
}

// The serverErrorResponse() method is used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses
// the errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowedResponse() method will be used to send a 405 method Not allowed
// status code and JSON response to the client
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
