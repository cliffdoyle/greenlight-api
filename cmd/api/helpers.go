package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// When httprouter is parsing a request, any interpolated URL parameters will be
// stored in the request context. We can use the ParamsFromContext() function to
// retrieve a slice containing these parameter names and values
// We can then use the ByName() method to get the value of the "id" parameter from
// the slice. In our project all movies will have a unique positive integer ID, but
// the value returned by ByName() is always a string. So we try to convert it to a
// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
// or is less than 1, we know the ID is invalid so we use the http.NotFound()
// function to return a 404 Not Found response
// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an erro
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJson(data any, w http.ResponseWriter, r *http.Request) {
	js, err := json.MarshalIndent(data,"","\t")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	js=append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
