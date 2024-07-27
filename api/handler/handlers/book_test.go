package handlers

import (
	_ "database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your_project/models" //
)

func TestGetBooks(t *testing.T) {
	// Создаем mock DB и handler
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "author", "publication_year", "isbn"}).
		AddRow(1, "Book One", "Author One", 2001, "111-1111111111").
		AddRow(2, "Book Two", "Author Two", 2002, "222-2222222222")

	// Настройка mock ожиданий
	mock.ExpectQuery("SELECT id, title, author, publication_year, isbn FROM books").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, author, publication_year, isbn FROM books")
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
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"ID":1,"Title":"Book One","Author":"Author One","PublicationYear":2001,"ISBN":"111-1111111111"},{"ID":2,"Title":"Book Two","Author":"Author Two","PublicationYear":2002,"ISBN":"222-2222222222"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
