package services

import (
	"project-go/internal/database"
)

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	BirthYear int    `json:"birthYear"`
	CreatedAt string `json:"createdAt"`
}

type AuthorService struct {
	db *database.DB
}

func NewAuthorService(db *database.DB) *AuthorService {
	return &AuthorService{db: db}
}

func (s *AuthorService) CreateAuthor(author *Author) error {
	query := "INSERT INTO authors (name, bio, birth_year) VALUES (?, ?, ?)"
	id, err := s.db.Insert(query, author.Name, author.Bio, author.BirthYear)

	if err != nil {
		return err
	}
	author.ID = int(id)
	return nil
}

func (s *AuthorService) GetAuthors() ([]Author, error) {
	query := "SELECT id, name, bio, birth_year, created_at FROM authors"
	rows, err := s.db.Query(s.db.Rebind(query))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var a Author
		err := rows.Scan(&a.ID, &a.Name, &a.Bio, &a.BirthYear, &a.CreatedAt)

		if err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	return authors, nil
}

func (s *AuthorService) GetAuthorByID(id int) (*Author, error) {
	query := "SELECT id, name, bio, birth_year, created_at FROM authors WHERE id = ?"
	row := s.db.QueryRow(s.db.Rebind(query), id)
	var a Author
	err := row.Scan(&a.ID, &a.Name, &a.Bio, &a.BirthYear, &a.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AuthorService) UpdateAuthor(author *Author) error {
	query := "UPDATE authors SET name = ?, bio = ?, birth_year = ? WHERE id = ?"
	_, err := s.db.Exec(s.db.Rebind(query), author.Name, author.Bio, author.BirthYear, author.ID)
	return err
}

func (s *AuthorService) DeleteAuthor(id int) error {
	query := "DELETE FROM authors WHERE id = ?"
	_, err := s.db.Exec(s.db.Rebind(query), id)
	return err
}

func (s *AuthorService) GetBooksByAuthor(authorID int) ([]Book, error) {
	query := "SELECT id, title, author, author_id, genre, pages, publicationyear FROM books WHERE author_id = ?"
	rows, err := s.db.Query(s.db.Rebind(query), authorID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.AuthorID, &b.Genre, &b.Pages, &b.PublicationYear)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}
