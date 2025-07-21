package data

import "time"

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
