package web

import (
	"fmt"
	"time"

	"project-go/internal/services"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %d error(s)", len(v.Errors))
}

func (v *ValidationErrors) Add(field, message string) {
	v.Errors = append(v.Errors, ValidationError{Field: field, Message: message})
}

func (v *ValidationErrors) IsEmpty() bool {
	return len(v.Errors) == 0
}

func ValidateBook(book *services.Book) *ValidationErrors {
	ve := &ValidationErrors{}

	if book.Title == "" {
		ve.Add("title", "title is required")
	}

	if book.Author == "" {
		ve.Add("author", "author is required")
	}

	if book.Genre == "" {
		ve.Add("genre", "genre is required")
	}

	if book.Pages <= 0 {
		ve.Add("pages", "pages must be greater than 0")
	}

	currentYear := time.Now().Year()
	if book.PublicationYear < 1000 || book.PublicationYear > currentYear+1 {
		ve.Add("publicationYear", fmt.Sprintf("publicationYear must be between 1000 and %d", currentYear+1))
	}

	return ve
}

func ValidateAuthor(author *services.Author) *ValidationErrors {
	ve := &ValidationErrors{}

	if author.Name == "" {
		ve.Add("name", "name is required")
	}

	if author.BirthYear != 0 && (author.BirthYear < 1000 || author.BirthYear > time.Now().Year()) {
		ve.Add("birthYear", "birthYear must be a valid year")
	}

	return ve
}

func ValidateLoan(loan *services.Loan) *ValidationErrors {
	ve := &ValidationErrors{}

	if loan.BookID <= 0 {
		ve.Add("bookId", "bookId must be greater than 0")
	}

	if loan.BorrowerName == "" {
		ve.Add("borrowerName", "borrowerName is required")
	}

	if loan.DueDate.IsZero() {
		ve.Add("dueDate", "dueDate is required")
	} else if loan.DueDate.Before(time.Now().AddDate(0, 0, -1)) {
		ve.Add("dueDate", "dueDate must be in the future")
	}

	return ve
}
