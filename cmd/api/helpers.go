package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1,048,576
	// bytes (1MB).
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body into the target destination.
	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		// Add a new maxBytesError variable.
		var maxBytesError *http.MaxBytesError
		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
			// A json.InvalidUnmarshalError error will be returned if we pass something
			// that is not a non-nil pointer to Decode(). We catch this and panic,
			// rather than returning an error to our handler.

			// Use the errors.As() function to check whether the error has the type
			// *http.MaxBytesError. If it does, then it means the request body exceeded our
			// size limit of 1MB and we return a clear error message.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// For anything else, return the error message as-is.
		default:
			return err
		}
	}
	return nil
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
