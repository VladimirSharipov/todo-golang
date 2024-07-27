package models

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
	ISBN            string `json:"isbn"`
}
