package web

import (
	"database/sql"
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

// GetBooks godoc
// @Summary      List books
// @Description  Get all books with pagination and filtering
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page            query  int    false  "Page number"      default(1)
// @Param        pageSize        query  int    false  "Items per page"   default(10)
// @Param        genre           query  string false  "Filter by genre"
// @Param        authorId        query  int    false  "Filter by author ID"
// @Param        publicationYear query  int    false  "Filter by year"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /books [get]
func (h *BookHandlers) GetBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := services.BookFilter{
		Genre:           query.Get("genre"),
		Page:            parseInt(query.Get("page")),
		PageSize:        parseInt(query.Get("pageSize")),
		PublicationYear: parseInt(query.Get("publicationYear")),
	}
	if aid := parseInt(query.Get("authorId")); aid > 0 {
		filter.AuthorID = aid
	}

	books, total, err := h.services.GetBooksFiltered(filter)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve books")
		return
	}

	page := filter.Page

	if page < 1 {
		page = 1
	}

	pageSize := filter.PageSize

	if pageSize < 1 {
		pageSize = 10
	}
	RespondWithPaginatedJSON(w, http.StatusOK, books, total, page, pageSize)
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}

// CreateBook godoc
// @Summary      Create a book
// @Description  Add a new book to the library
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      services.Book  true  "Book data"
// @Success      201   {object}  services.Book
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /books [post]
func (h *BookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book services.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if ve := ValidateBook(&book); !ve.IsEmpty() {
		RespondWithJSON(w, http.StatusBadRequest, ve)
		return
	}

	if err := h.services.CreateBook(&book); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to create book")
		return
	}
	RespondWithJSON(w, http.StatusCreated, book)
}

// GetBookByID godoc
// @Summary      Get a book
// @Description  Get a book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  services.Book
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /books/{id} [get]
func (h *BookHandlers) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid book ID")
		return
	}

	book, err := h.services.GetBookByID(id)
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "not_found", "Book not found")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve book")
		return
	}

	RespondWithJSON(w, http.StatusOK, book)
}

func (h *BookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid book ID")
		return
	}

	var book services.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if ve := ValidateBook(&book); !ve.IsEmpty() {
		RespondWithJSON(w, http.StatusBadRequest, ve)
		return
	}

	book.ID = id
	if err := h.services.UpdateBook(&book); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to update book")
		return
	}
	RespondWithJSON(w, http.StatusOK, book)
}

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Remove a book from the library
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /books/{id} [delete]
func (h *BookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid book ID")
		return
	}

	if err := h.services.DeleteBook(id); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to delete book")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
