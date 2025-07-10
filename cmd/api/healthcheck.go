package main

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

//The important thing to point out is that healthcheckHandler is implemented
//as a method on the application struct
//This is an effective and idiomatic way to make dependencies available to our handlers
//without resorting to global variables or closures.
//any dependencies that the healthcheckHandler needs can simply be included
//as a field in the application struct when initialized in the main()

//We can see this pattern already being used in the code above, where the operating
//environment name is retrieved from the application struct by calling app.config.env
