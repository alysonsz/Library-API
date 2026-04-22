package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"project-go/internal/database"
)

func TestCreateBook(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	book := &Book{
		Title:           "Test Book",
		Author:          "Test Author",
		AuthorID:        1,
		Genre:           "Fiction",
		Pages:           300,
		PublicationYear: 2024,
	}

	mock.ExpectExec("INSERT INTO books").
		WithArgs(book.Title, book.Author, book.AuthorID, book.Genre, book.Pages, book.PublicationYear).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := service.CreateBook(book); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if book.ID != 1 {
		t.Errorf("expected ID 1, got %d", book.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetBooks(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "author_id", "genre", "pages", "publicationyear"}).
		AddRow(1, "Book One", "Author A", 1, "Fiction", 200, 2020).
		AddRow(2, "Book Two", "Author B", 2, "Non-Fiction", 300, 2021)

	mock.ExpectQuery("SELECT (.+) FROM books").WillReturnRows(rows)

	books, err := service.GetBooks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}

	if books[0].Title != "Book One" {
		t.Errorf("expected 'Book One', got %s", books[0].Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetBookByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "author_id", "genre", "pages", "publicationyear"}).
		AddRow(1, "Book One", "Author A", 1, "Fiction", 200, 2020)

	mock.ExpectQuery("SELECT (.+) FROM books WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	book, err := service.GetBookByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if book == nil {
		t.Fatal("expected book, got nil")
	}

	if book.Title != "Book One" {
		t.Errorf("expected 'Book One', got %s", book.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUpdateBook(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	book := &Book{
		ID:              1,
		Title:           "Updated Book",
		Author:          "Author A",
		AuthorID:        1,
		Genre:           "Fiction",
		Pages:           250,
		PublicationYear: 2023,
	}

	mock.ExpectExec("UPDATE books SET").
		WithArgs(book.Title, book.Author, book.AuthorID, book.Genre, book.Pages, book.PublicationYear, book.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := service.UpdateBook(book); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteBook(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	mock.ExpectExec("DELETE FROM books WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := service.DeleteBook(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestSearchBooksByName(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "author_id", "genre", "pages", "publicationyear"}).
		AddRow(1, "Go in Action", "Author A", 1, "Tech", 400, 2022)

	mock.ExpectQuery("SELECT (.+) FROM books WHERE title LIKE ?").
		WithArgs("%Go%").
		WillReturnRows(rows)

	books, err := service.SearchBooksByName("Go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetBooksFiltered(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer mockDB.Close()

	db := database.NewDB(mockDB, "sqlite3")
	service := NewBookService(db)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM books WHERE 1=1 AND genre = \\?").
		WithArgs("Fiction").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	rows := sqlmock.NewRows([]string{"id", "title", "author", "author_id", "genre", "pages", "publicationyear"}).
		AddRow(1, "Book A", "Author A", 1, "Fiction", 200, 2020).
		AddRow(2, "Book B", "Author B", 2, "Fiction", 300, 2021)

	mock.ExpectQuery("SELECT id, title, author, author_id, genre, pages, publicationyear FROM books WHERE 1=1 AND genre = \\? LIMIT \\? OFFSET \\?").
		WithArgs("Fiction", 10, 0).
		WillReturnRows(rows)

	filter := BookFilter{
		Genre:    "Fiction",
		Page:     1,
		PageSize: 10,
	}

	books, total, err := service.GetBooksFiltered(filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(books) != 2 {
		t.Errorf("expected 2 books, got %d", len(books))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
