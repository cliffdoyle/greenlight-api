package main

import (
	"fmt"
	"net/http"
)

	// Declare a handler which writes a plain-text response with information about the
	// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a fixed-format JSON response from a string. Notice how we're using a raw
	// string literal (enclosed with backticks) so that we can include double-quote
	// characters in the JSON without needing to escape them? We also use the %q verb to
	// wrap the interpolated values in double-quotes
	js:=`{"status":"available", "environment":%q,"version":%q}`
	js=fmt.Sprintf(js,app.config.env,version)

	// Set the "Content-Type: application/json" header on the response. If you forget to
	// this, Go will default to sending a "Content-Type: text/plain; charset=utf-8"
	// header instead
	w.Header().Set("Content-Type","application/json")

	//Write the JSON as the HTTP response body
	w.Write([]byte(js))
}

//The important thing to point out is that healthcheckHandler is implemented
//as a method on the application struct
//This is an effective and idiomatic way to make dependencies available to our handlers
//without resorting to global variables or closures.
//any dependencies that the healthcheckHandler needs can simply be included
//as a field in the application struct when initialized in the main()

//We can see this pattern already being used in the code above, where the operating
//environment name is retrieved from the application struct by calling app.config.env
