package services

import "database/sql"

type Book struct {
	ID              int
	Title           string
	Author          string
	Pages           int
	Genre           string
	PublicationYear int
}

func (bookInformation Book) GetTitleAndAuthor() string {

	return bookInformation.Title + " by " + bookInformation.Author

}

type BookService struct {
	database *sql.DB
}

func (Service *BookService) CreateBook(book *Book) error {

	query := "Insert into books (title, author, pages, genre, publicationyear) values(?, ?, ?, ?, ?)"
	result, error := Service.database.Exec(query, book.Title, book.Author, book.Pages, book.Genre, book.PublicationYear) //Executa o banco de dados inserindo os valores para serem substituídos na "?"
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

func (Service *BookService) GetBook() ([]Book, error) {

	query := "Select id, title, author, pages, genre, publicationyear from books"
	rows, error := Service.database.Query(query)
	if error != nil {
		return nil, error
	}

	var books []Book
	for rows.Next() {
		var book Book
		error := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Pages, &book.Genre, &book.PublicationYear)
		if error != nil {
			return nil, error
		}
		books = append(books, book)
	}
	return books, nil
}

func (Service *BookService) GetBookID(id int) (*Book, error) {

	query := "Select id, title, author, pages, genre, publicationyear from books where id = ?"
	row := Service.database.QueryRow(query, id)

	var book Book
	error := row.Scan(&book.ID, &book.Title, &book.Author, &book.Pages, &book.Genre, &book.PublicationYear)
	if error != nil {
		return nil, error
	}
	return &book, nil

}

func (Service *BookService) UpdateBook(book *Book) error {

	query := "Update books set title = ?, author = ?, pages = ?, genre = ?, publicationyear = ? where id = ?"
	_, error := Service.database.Exec(query, book.Title, book.Author, book.Pages, book.Genre, book.PublicationYear, book.ID)
	if error != nil {
		return error
	}
	return nil
}

func (Service *BookService) DeleteBook(id int) error {

	query := "Delete from books where id = ?"
	_, error := Service.database.Exec(query, id)
	if error != nil {
		return error
	}
	return nil

}
