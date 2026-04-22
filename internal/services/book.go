package services

import (
	"fmt"
	"project-go/internal/database"
	"time"
)

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	AuthorID        int    `json:"authorId"`
	Genre           string `json:"genre"`
	Pages           int    `json:"pages"`
	PublicationYear int    `json:"publicationYear"`
}

type BookFilter struct {
	Genre           string
	AuthorID        int
	PublicationYear int
	Page            int
	PageSize        int
}

type BookService struct {
	database *database.DB
}

func NewBookService(database *database.DB) *BookService {
	return &BookService{database: database}
}

func (Service *BookService) CreateBook(book *Book) error {
	query := "INSERT INTO books (title, author, author_id, genre, pages, publicationyear) VALUES (?, ?, ?, ?, ?, ?)"
	id, err := Service.database.Insert(query, book.Title, book.Author, book.AuthorID, book.Genre, book.Pages, book.PublicationYear)

	if err != nil {
		return err
	}
	book.ID = int(id)
	return nil
}

func (Service *BookService) GetBooks() ([]Book, error) {
	query := "SELECT id, title, author, author_id, genre, pages, publicationyear FROM books"
	rows, err := Service.database.Query(Service.database.Rebind(query))

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.AuthorID, &book.Genre, &book.Pages, &book.PublicationYear)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (Service *BookService) GetBookByID(ID int) (*Book, error) {
	query := "SELECT id, title, author, author_id, genre, pages, publicationyear FROM books WHERE id = ?"
	row := Service.database.QueryRow(Service.database.Rebind(query), ID)
	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.AuthorID, &book.Genre, &book.Pages, &book.PublicationYear)

	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (Service *BookService) UpdateBook(book *Book) error {
	query := "UPDATE books SET title = ?, author = ?, author_id = ?, genre = ?, pages = ?, publicationyear = ? WHERE id = ?"
	_, err := Service.database.Exec(Service.database.Rebind(query), book.Title, book.Author, book.AuthorID, book.Genre, book.Pages, book.PublicationYear, book.ID)
	return err
}

func (Service *BookService) DeleteBook(ID int) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := Service.database.Exec(Service.database.Rebind(query), ID)
	return err
}

func (Service *BookService) SearchBooksByName(nameBook string) ([]Book, error) {
	query := "SELECT id, title, author, author_id, genre, pages, publicationyear FROM books WHERE title LIKE ?"
	rows, err := Service.database.Query(Service.database.Rebind(query), "%"+nameBook+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.AuthorID, &book.Genre, &book.Pages, &book.PublicationYear)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (Service *BookService) GetBooksFiltered(filter BookFilter) ([]Book, int, error) {
	query := "SELECT id, title, author, author_id, genre, pages, publicationyear FROM books WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM books WHERE 1=1"
	args := []interface{}{}

	if filter.Genre != "" {
		query += " AND genre = ?"
		countQuery += " AND genre = ?"
		args = append(args, filter.Genre)
	}

	if filter.AuthorID > 0 {
		query += " AND author_id = ?"
		countQuery += " AND author_id = ?"
		args = append(args, filter.AuthorID)
	}

	if filter.PublicationYear > 0 {
		query += " AND publicationyear = ?"
		countQuery += " AND publicationyear = ?"
		args = append(args, filter.PublicationYear)
	}

	var total int
	err := Service.database.QueryRow(Service.database.Rebind(countQuery), args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := Service.database.Query(Service.database.Rebind(query), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.AuthorID, &book.Genre, &book.Pages, &book.PublicationYear)
		if err != nil {
			return nil, 0, err
		}
		books = append(books, book)
	}
	return books, total, nil
}

func (Service *BookService) SimulateReading(BookID int, duration time.Duration, results chan<- string) {
	book, err := Service.GetBookByID(BookID)
	if err != nil || book == nil {
		results <- fmt.Sprintf("Book %d not found", BookID)
		return
	}
	time.Sleep(duration)
	results <- fmt.Sprintf("Book %s readed", book.Title)
}

func (Service *BookService) SimulateMultipleReading(BookIDs []int, duration time.Duration) []string {
	results := make(chan string, len(BookIDs))
	for _, ID := range BookIDs {
		go func(BookID int) {
			Service.SimulateReading(BookID, duration, results)
		}(ID)
	}

	var responses []string
	for range BookIDs {
		responses = append(responses, <-results)
	}
	close(results)
	return responses
}
