package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/greenlight-api/validator"
	"github.com/lib/pq"
)

//To control the visibility of individual struct fields in the JSON
//use omitzero and - struct tag directives.

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   int32     `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

// MovieModel struct type that wraps a sql.DB connection pool
type MovieModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table
func (m MovieModel) Insert(movie *Movie) error {
	//SQL querry for insrting a new record in the movies table and returning
	//system-generated data
	query := `
	INSERT INTO movies (title,year,runtime,genres,version)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, created_at,version
	`

	//Args slice containing the values for the placeholder parameters from
	//the movie struct.Declaring it immediately next to the sql query
	//makes it clear what values are used in the query
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.Version}

	//use the QueryRow() method to execute the SQL query on the connection pool,
	//passing in the args slice as a variadic parameter and scanning the system-generated
	//id , created_at and version values into the movie struct
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	//Because the Insert() method takes a *Movie pointer as the parameter, when we call Scan() to read in the
	//system-generated data we're updating the values at the location the parameter points to.
	//Essentially, the Insert() method mutates the Movie struct that is passed to it and adds the system-generated
	//values to it
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	//The PostgreSQL bigserial type we use for the  movie ID starts auto-incrementing
	//at 1 by default, so no movie will have ID values less than that.
	//To avoid making unnecessary database call, we take a shortcut and return an
	//ErrRecordNotFound error straight away
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//Define the SQL query for retrieving the movie data
	Query := `SELECT id, created_at, title, year,runtime,genres, version
			FROM movies
			WHERE id =$1		
	`

	//Declare a movie struct to hold the data returned by the query
	var movie Movie
	//Execute the query using the QueryRow() method, passing in the provided id value
	//as a placeholder parameter, and scan the response data into the fields of the movie struct
	//we convert the scan target for the genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRow(Query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	//Handle any errors.If there was no matching movie found, scan() will
	//return a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	//error instead

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	//Otherwise, return a pointer to the movie struct
	return &movie, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	//Declare the SQL query for updating the record and returning the new version
	//number

	query := `
	UPDATE movies
	SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1
	WHERE id = $5
	RETURNING version
	`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	//Use the QueryRow() method to execute the query, passing in the args slice as 
	//variadic parameter and scanning the new version value into the movie struct
	return  m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return nil
}
