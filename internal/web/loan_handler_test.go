package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"project-go/internal/database"
	"project-go/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupLoanTestHandler() (*LoanHandlers, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	db := database.NewDB(mockDB, "sqlite3")
	service := services.NewLoanService(db)
	handler := NewLoanHandlers(service)
	return handler, mock, func() { mockDB.Close() }
}

func TestGetLoans(t *testing.T) {
	handler, mock, cleanup := setupLoanTestHandler()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "book_id", "borrower_name", "loan_date", "due_date", "return_date", "status"}).
		AddRow(1, 1, "Alice", time.Now(), time.Now().Add(7*24*time.Hour), nil, "active")

	mock.ExpectQuery("SELECT (.+) FROM loans").WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/loans", nil)
	rr := httptest.NewRecorder()

	handler.GetLoans(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var loans []services.Loan
	if err := json.Unmarshal(rr.Body.Bytes(), &loans); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(loans) != 1 {
		t.Errorf("expected 1 loan, got %d", len(loans))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateLoan(t *testing.T) {
	handler, mock, cleanup := setupLoanTestHandler()
	defer cleanup()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM loans WHERE book_id = \\? AND status = 'active'").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec("INSERT INTO loans").
		WithArgs(1, "Alice", sqlmock.AnyArg(), "active").
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := map[string]interface{}{
		"bookId":       1,
		"borrowerName": "Alice",
		"dueDate":      time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateLoan(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var loan services.Loan
	if err := json.Unmarshal(rr.Body.Bytes(), &loan); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if loan.ID != 1 {
		t.Errorf("expected ID 1, got %d", loan.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateLoanAlreadyLoaned(t *testing.T) {
	handler, mock, cleanup := setupLoanTestHandler()
	defer cleanup()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM loans WHERE book_id = \\? AND status = 'active'").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	body := map[string]interface{}{
		"bookId":       1,
		"borrowerName": "Alice",
		"dueDate":      time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateLoan(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestReturnLoan(t *testing.T) {
	handler, mock, cleanup := setupLoanTestHandler()
	defer cleanup()

	mock.ExpectExec("UPDATE loans SET status = 'returned'").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodPost, "/loans/1/return", nil)
	req.SetPathValue("id", "1")
	rr := httptest.NewRecorder()

	handler.ReturnLoan(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestReturnLoanNotFound(t *testing.T) {
	handler, mock, cleanup := setupLoanTestHandler()
	defer cleanup()

	mock.ExpectExec("UPDATE loans SET status = 'returned'").
		WithArgs(sqlmock.AnyArg(), 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest(http.MethodPost, "/loans/999/return", nil)
	req.SetPathValue("id", "999")
	rr := httptest.NewRecorder()

	handler.ReturnLoan(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetOverdueLoans(t *testing.T) {
	handler, mock, cleanup := setupLoanTestHandler()
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "book_id", "borrower_name", "loan_date", "due_date", "return_date", "status"}).
		AddRow(1, 1, "Alice", time.Now().Add(-14*24*time.Hour), time.Now().Add(-7*24*time.Hour), nil, "active")

	mock.ExpectQuery("SELECT (.+) FROM loans WHERE due_date < \\? AND status = 'active'").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/loans/overdue", nil)
	rr := httptest.NewRecorder()

	handler.GetOverdueLoans(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var loans []services.Loan
	if err := json.Unmarshal(rr.Body.Bytes(), &loans); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(loans) != 1 {
		t.Errorf("expected 1 overdue loan, got %d", len(loans))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
