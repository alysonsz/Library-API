package web

import (
	"encoding/json"
	"net/http"
	"project-go/internal/services"
)

type BookHandlers struct {
	Services *services.BookService
}

func (Handler *BookHandlers) GetBooks(writter http.ResponseWriter, request *http.Request) {

	books, error := Handler.Services.GetBooks()
	if error != nil {

		http.Error(writter, "Error to get books", http.StatusInternalServerError)
		return

	}
	writter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writter).Encode(books)

}

func (Handler *BookHandlers) CreateBook(writter http.ResponseWriter, request *http.Request) {

	var books services.Book
	error := json.NewDecoder(request.Body).Decode(&books)
	if error != nil {

		http.Error(writter, "Error to decode body", http.StatusBadRequest)
		return

	}
	error = Handler.Services.CreateBook(&books)
	if error != nil {

		http.Error(writter, "Error to create book", http.StatusInternalServerError)
		return

	}
	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(books)

}
