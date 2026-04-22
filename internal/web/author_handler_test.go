package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"project-go/internal/database"
	"project-go/internal/services"
)

func setupAuthorTestHandler() (*AuthorHandlers, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	db := database.NewDB(mockDB, "sqlite3")
	service := services.NewAuthorService(db)
	handler := NewAuthorHandlers(service)
	return handler, mock, func() { mockDB.Close() }
}

func TestGetAuthors(t *testing.T) {
	handler, mock, cleanup := setupAuthorTestHandler()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "bio", "birth_year", "created_at"}).
		AddRow(1, "Author A", "Bio A", 1900, "2024-01-01T00:00:00Z").
		AddRow(2, "Author B", "Bio B", 1950, "2024-01-02T00:00:00Z")
	mock.ExpectQuery("SELECT (.+) FROM authors").WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/authors", nil)
	rr := httptest.NewRecorder()

	handler.GetAuthors(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var authors []services.Author
	if err := json.Unmarshal(rr.Body.Bytes(), &authors); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(authors))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateAuthor(t *testing.T) {
	handler, mock, cleanup := setupAuthorTestHandler()
	defer cleanup()

	mock.ExpectExec("INSERT INTO authors").
		WithArgs("George Orwell", "English novelist", 1903).
		WillReturnResult(sqlmock.NewResult(1, 1))

	author := services.Author{Name: "George Orwell", Bio: "English novelist", BirthYear: 1903}
	body, _ := json.Marshal(author)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateAuthor(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var created services.Author
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateAuthorValidation(t *testing.T) {
	handler, _, cleanup := setupAuthorTestHandler()
	defer cleanup()

	author := services.Author{Name: ""}
	body, _ := json.Marshal(author)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateAuthor(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetAuthorByID(t *testing.T) {
	handler, mock, cleanup := setupAuthorTestHandler()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "bio", "birth_year", "created_at"}).
		AddRow(1, "Author A", "Bio A", 1900, "2024-01-01T00:00:00Z")

	mock.ExpectQuery("SELECT (.+) FROM authors WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/authors/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	handler.GetAuthorByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var author services.Author
	if err := json.Unmarshal(rr.Body.Bytes(), &author); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if author.Name != "Author A" {
		t.Errorf("expected 'Author A', got %s", author.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetAuthorByIDNotFound(t *testing.T) {
	handler, mock, cleanup := setupAuthorTestHandler()
	defer cleanup()

	mock.ExpectQuery("SELECT (.+) FROM authors WHERE id = \\?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/authors/999", nil)
	req.SetPathValue("id", "999")
	rr := httptest.NewRecorder()

	handler.GetAuthorByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteAuthor(t *testing.T) {
	handler, mock, cleanup := setupAuthorTestHandler()
	defer cleanup()

	mock.ExpectExec("DELETE FROM authors WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodDelete, "/authors/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	handler.DeleteAuthor(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
