package data

import (
	"database/sql"
	"time"
)

//To control the visibility of individual struct fields in the JSON
//use omitzero and - struct tag directives.

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"titledeed"`
	Year      int32     `json:"year,omitzero"`
	Runtime   int32     `json:"runtime,omitzero"`
	Genre     []string  `json:"genre,omitzero"`
	Version   int32     `json:"version"`
}

//MovieModel struct type that wraps a sql.DB connection pool
type MovieModel struct{
	DB *sql.DB
}

//Add a placeholder method for inserting a new record in the movies table
func (m MovieModel)Insert(movie *Movie)error{
	return  nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
return nil, nil
}
// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
return nil
}
// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
return nil
}

