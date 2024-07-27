package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VladimirSharipov/Todo-golang/internal/models"
)

func HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w, r)
	case http.MethodPost:
		createBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBook(w, r)
	case http.MethodPut:
		updateBook(w, r)
	case http.MethodDelete:
		deleteBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, author, publication_year, isbn FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear, &book.ISBN); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/books/"):]
	var book models.Book
	err := db.DB.QueryRow("SELECT id, title, author, publication_year, isbn FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.PublicationYear, &book.ISBN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow("INSERT INTO books(title, author, publication_year, isbn) VALUES($1, $2, $3, $4) RETURNING id", book.Title, book.Author, book.PublicationYear, book.ISBN).Scan(&book.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/books/"):]
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE books SET title=$1, author=$2, publication_year=$3, isbn=$4 WHERE id=$5", book.Title, book.Author, book.PublicationYear, book.ISBN, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID, _ = strconv.Atoi(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/books/"):]
	_, err := db.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
