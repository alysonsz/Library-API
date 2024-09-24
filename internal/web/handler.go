package web

import (
	"encoding/json"
	"net/http"
	"project-go/internal/services"
	"strconv"
)

type BookHandlers struct {
	services *services.BookService
}

func NewBookHandlers(services *services.BookService) *BookHandlers {

	return &BookHandlers{services: services}

}

func (Handler *BookHandlers) GetBooks(writer http.ResponseWriter, request *http.Request) {

	books, error := Handler.services.GetBooks()
	if error != nil {

		http.Error(writer, "Error to get books", http.StatusInternalServerError)
		return

	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(books)

}

func (Handler *BookHandlers) CreateBook(writer http.ResponseWriter, request *http.Request) {

	var books services.Book
	error := json.NewDecoder(request.Body).Decode(&books)
	if error != nil {

		http.Error(writer, "Error to decode body", http.StatusBadRequest)
		return

	}
	error = Handler.services.CreateBook(&books)
	if error != nil {

		http.Error(writer, "Error to create book", http.StatusInternalServerError)
		return

	}
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(books)

}
func (Handler *BookHandlers) GetBookByID(writer http.ResponseWriter, request *http.Request) {

	idString := request.PathValue("id")
	id, error := strconv.Atoi(idString)
	if error != nil {

		http.Error(writer, "invalid book ID", http.StatusBadRequest)
		return

	}
	books, error := Handler.services.GetBookByID(id)
	if error != nil {

		http.Error(writer, "failed to get book", http.StatusBadRequest)
		return

	}
	if error == nil {

		http.Error(writer, "book not found", http.StatusBadRequest)
		return

	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(books)

}

func (Handler *BookHandlers) UpdateBook(writer http.ResponseWriter, request *http.Request) {
	idString := request.PathValue("id")
	id, error := strconv.Atoi(idString)
	if error != nil {

		http.Error(writer, "invalid book ID", http.StatusBadRequest)

	}

	if error != nil {

		var books services.Book
		error := json.NewDecoder(request.Body).Decode(&books)
		if error != nil {

			http.Error(writer, "invalid request", http.StatusBadRequest)
			return

		}
		books.ID = id
		error = Handler.services.UpdateBook(&books)
		if error != nil {

			http.Error(writer, "failed to update book", http.StatusBadRequest)
			return

		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(books)

	}
}
func (Handler *BookHandlers) DeleteBook(writer http.ResponseWriter, request *http.Request) {

	idString := request.PathValue("id")
	id, error := strconv.Atoi(idString)
	if error != nil {

		http.Error(writer, "invalid book id", http.StatusBadRequest)
		return

	}

	error = Handler.services.DeleteBook(id)
	if error != nil {

		http.Error(writer, "failed to delete book", http.StatusBadRequest)
		return

	}
	writer.WriteHeader(http.StatusNoContent)

}
