package web

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"project-go/internal/services"
	"strconv"
)

type AuthorHandlers struct {
	service *services.AuthorService
}

func NewAuthorHandlers(service *services.AuthorService) *AuthorHandlers {
	return &AuthorHandlers{service: service}
}

func (h *AuthorHandlers) GetAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := h.service.GetAuthors()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve authors")
		return
	}
	RespondWithJSON(w, http.StatusOK, authors)
}

// CreateAuthor godoc
// @Summary      Create an author
// @Description  Add a new author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        author  body      services.Author  true  "Author data"
// @Success      201     {object}  services.Author
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /authors [post]
func (h *AuthorHandlers) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author services.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if ve := ValidateAuthor(&author); !ve.IsEmpty() {
		RespondWithJSON(w, http.StatusBadRequest, ve)
		return
	}

	if err := h.service.CreateAuthor(&author); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to create author")
		return
	}
	RespondWithJSON(w, http.StatusCreated, author)
}

// GetAuthorByID godoc
// @Summary      Get an author
// @Description  Get author by ID
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Author ID"
// @Success      200  {object}  services.Author
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /authors/{id} [get]
func (h *AuthorHandlers) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid author ID")
		return
	}
	author, err := h.service.GetAuthorByID(id)
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "not_found", "Author not found")
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve author")
		return
	}
	RespondWithJSON(w, http.StatusOK, author)
}

// UpdateAuthor godoc
// @Summary      Update an author
// @Description  Update existing author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id      path      int              true  "Author ID"
// @Param        author  body      services.Author  true  "Author data"
// @Success      200     {object}  services.Author
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /authors/{id} [put]
func (h *AuthorHandlers) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid author ID")
		return
	}
	var author services.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}
	if ve := ValidateAuthor(&author); !ve.IsEmpty() {
		RespondWithJSON(w, http.StatusBadRequest, ve)
		return
	}
	author.ID = id
	if err := h.service.UpdateAuthor(&author); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to update author")
		return
	}
	RespondWithJSON(w, http.StatusOK, author)
}

// DeleteAuthor godoc
// @Summary      Delete an author
// @Description  Remove an author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Author ID"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /authors/{id} [delete]
func (h *AuthorHandlers) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid author ID")
		return
	}
	if err := h.service.DeleteAuthor(id); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to delete author")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetBooksByAuthor godoc
// @Summary      List books by author
// @Description  Get all books written by an author
// @Tags         authors
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Author ID"
// @Success      200  {array}   services.Book
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /authors/{id}/books [get]
func (h *AuthorHandlers) GetBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_id", "Invalid author ID")
		return
	}
	books, err := h.service.GetBooksByAuthor(id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve books")
		return
	}
	RespondWithJSON(w, http.StatusOK, books)
}
