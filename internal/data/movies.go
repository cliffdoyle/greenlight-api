package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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
	//SQL querry for insrting a new record in the movies table and returning 
	//system-generated data
	query:=`
	INSERT INTO movies (title,year,runtime,genres)
	VALUES ($1,$2,$3,$4)
	RETURNING id, created_at,version
	`

	//Args slice containing the values for the placeholder parameters from 
	//the movie struct.Declaring it immediately next to the sql query
	//makes it clear what values are used in the query
	args:=[]any{movie.Title,movie.Year,movie.Runtime,pq.Array(movie.Genre)}

	//use the QueryRow() method to execute the SQL query on the connection pool,
	//passing in the args slice as a variadic parameter and scanning the system-generated
	//id , created_at and version values into the movie struct
	return  m.DB.QueryRow(query,args...).Scan(&movie.ID,&movie.CreatedAt,&movie.Version)

	//Because the Insert() method takes a *Movie pointer as the parameter, when we call Scan() to read in the 
	//system-generated data we're updating the values at the location the parameter points to.
	//Essentially, the Insert() method mutates the Movie struct that is passed to it and adds the system-generated
	//values to it
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

