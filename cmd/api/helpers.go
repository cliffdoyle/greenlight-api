package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	//Import the pq driver so that it can register itself with the database/sql
	_ "github.com/lib/pq"
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

//For sending JSON responses to the client

func (app *application) writeJson(data any, w http.ResponseWriter, r *http.Request) {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// The openDB() function returns a sql.DB connection pool
// Uses db.PingContext() to actually create a connection and verify
// that everything is set up correctly
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	//Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Use pingContext() to establish a new connection to the database, passing in the
	//context we created above as a parameter.If the connection couldn't be
	//established successfully within the 5 second deadline, then this returns an error
	//If we get this error, or any other , we close the connection pool and return the error
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	//Return the sql.DB connection pool
	return db, nil
}
