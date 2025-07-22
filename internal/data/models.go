package data

import (
	"database/sql"
	"errors"
)

//Define a custom ErrRecordNotFound error. will be returned from Get()
//when looking up a movie that doesn't exist in the database
var (
	ErrRecordNotFound=errors.New("record not found")
)

//Create a  models struct which wraps the MovieModel
//we can now add other models to this later
type Models struct{
	Movies MovieModel
}

//For ease of use, we also add a New() method which returns a models struct
//containing the initialized MovieModel
func NewModels(db *sql.DB)Models{
	return  Models{Movies: MovieModel{DB: db}}
}