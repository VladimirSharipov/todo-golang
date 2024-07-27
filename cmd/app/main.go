package main

import (
	db "Todo-golang/internal/database"
	"log"
	"net/http"

	"github.com/VladimirSharipov/Todo-golang/handler/handlers"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	http.HandleFunc("/students", handlers.HandleStudents)
	http.HandleFunc("/students/", handlers.HandleStudent)
	http.HandleFunc("/books", handlers.HandleBooks)
	http.HandleFunc("/books/", handlers.HandleBook)
	http.HandleFunc("/products", handlers.HandleProducts)
	http.HandleFunc("/products/", handlers.HandleProduct)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
