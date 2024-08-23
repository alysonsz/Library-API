package services

import "database/sql"

type Book struct {
	ID              int
	Title           string
	Author          string
	Genre           string
	Pages           int
	publicationYear int
}

type BookService struct {
	database *sql.DB
}

func (Service *BookService) CreateBook(book *Book) error {

	query := "Insert into books (title, author, genre, pages, publicationyear) values (?, ?, ?, ?, ?)"
	result, error := Service.database.Exec(query, book.Title, book.Author, book.Genre, book.Pages, book.publicationYear)
	if error != nil {

		return error

	}
	lastInsertID, error := result.LastInsertId()
	if error != nil {

		return error

	}
	book.ID = int(lastInsertID)
	return nil

}

func (Service *BookService) GetBooks() ([]Book, error) {

	query := "Select id, title, author, genre, pages, publicationyear from books"
	rows, error := Service.database.Query(query)
	if error != nil {

		return nil, error

	}
	var books []Book
	for rows.Next() {
		var book Book
		error := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Pages, &book.publicationYear)
		if error != nil {

			return nil, error

		}
		books = append(books, book)
	}
	return books, nil

}

func (Service *BookService) GetBookByID(ID int) (*Book, error) {

	query := "Select id, title, author, genre, pages, publicationyear from books where id = ?"
	row := Service.database.QueryRow(query, ID)
	var book Book
	error := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre, &book.Pages, &book.publicationYear)
	if error != nil {

		return nil, error

	}
	return &book, nil

}

func (Service *BookService) UpdateBook(book *Book) error {

	query := "Update books set title = ?, author = ?, genre = ?, pages = ?, publicationyear = ? where id = ?"
	_, error := Service.database.Exec(query, book.Title, book.Author, book.Genre, book.Pages, book.publicationYear, book.ID)
	return error

}

func (Service *BookService) DeleteBook(ID int) error {

	query := "Delete from books where id = ?"
	_, error := Service.database.Exec(query, ID)
	return error

}
