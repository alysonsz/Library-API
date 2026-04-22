package web

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"project-go/internal/database"
	"project-go/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupTestHandler() (*BookHandlers, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	db := database.NewDB(mockDB, "sqlite3")
	service := services.NewBookService(db)
	handler := NewBookHandlers(service)
	return handler, mock, func() { mockDB.Close() }
}

func TestGetBooks(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM books WHERE 1=1").WillReturnRows(rows)

	bookRows := sqlmock.NewRows([]string{"id", "title", "author", "author_id", "genre", "pages", "publicationyear"}).
		AddRow(1, "Book A", "Author A", 1, "Fiction", 200, 2020).
		AddRow(2, "Book B", "Author B", 2, "Non-Fiction", 300, 2021)
	mock.ExpectQuery("SELECT (.+) FROM books WHERE 1=1 LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnRows(bookRows)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rr := httptest.NewRecorder()

	handler.GetBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp PaginatedResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateBook(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	mock.ExpectExec("INSERT INTO books").
		WithArgs("New Book", "Author", 0, "Fiction", 100, 2024).
		WillReturnResult(sqlmock.NewResult(1, 1))

	book := services.Book{
		Title:           "New Book",
		Author:          "Author",
		Genre:           "Fiction",
		Pages:           100,
		PublicationYear: 2024,
	}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateBook(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var created services.Book
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

func TestGetBookByID(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "title", "author", "author_id", "genre", "pages", "publicationyear"}).
		AddRow(1, "Book A", "Author A", 1, "Fiction", 200, 2020)

	mock.ExpectQuery("SELECT (.+) FROM books WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	handler.GetBookByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var book services.Book
	if err := json.Unmarshal(rr.Body.Bytes(), &book); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if book.Title != "Book A" {
		t.Errorf("expected 'Book A', got %s", book.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetBookByIDNotFound(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	mock.ExpectQuery("SELECT (.+) FROM books WHERE id = \\?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	req.SetPathValue("id", "999")
	rr := httptest.NewRecorder()

	handler.GetBookByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUpdateBook(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	mock.ExpectExec("UPDATE books SET").
		WithArgs("Updated", "Author", 0, "Fiction", 200, 2023, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	book := services.Book{
		Title:           "Updated",
		Author:          "Author",
		Genre:           "Fiction",
		Pages:           200,
		PublicationYear: 2023,
	}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPut, "/books/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.UpdateBook(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteBook(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	mock.ExpectExec("DELETE FROM books WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	handler.DeleteBook(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteBookInvalidID(t *testing.T) {
	handler, mock, cleanup := setupTestHandler()
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/books/abc", nil)
	req.SetPathValue("id", "abc")
	rr := httptest.NewRecorder()

	handler.DeleteBook(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
