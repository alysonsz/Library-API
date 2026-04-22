package services

import (
	"database/sql"
	"fmt"
	"project-go/internal/database"
	"time"
)

type Loan struct {
	ID           int       `json:"id"`
	BookID       int       `json:"bookId"`
	BorrowerName string    `json:"borrowerName"`
	LoanDate     time.Time `json:"loanDate"`
	DueDate      time.Time `json:"dueDate"`
	ReturnDate   time.Time `json:"returnDate,omitempty"`
	Status       string    `json:"status"`
}

type LoanService struct {
	db *database.DB
}

func NewLoanService(db *database.DB) *LoanService {
	return &LoanService{db: db}
}

func (s *LoanService) IsBookAvailable(bookID int) (bool, error) {
	query := "SELECT COUNT(*) FROM loans WHERE book_id = ? AND status = 'active'"
	var count int
	err := s.db.QueryRow(s.db.Rebind(query), bookID).Scan(&count)

	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *LoanService) CreateLoan(loan *Loan) error {
	available, err := s.IsBookAvailable(loan.BookID)

	if err != nil {
		return err
	}

	if !available {
		return fmt.Errorf("book %d is already loaned", loan.BookID)
	}

	query := "INSERT INTO loans (book_id, borrower_name, due_date, status) VALUES (?, ?, ?, ?)"
	id, err := s.db.Insert(query, loan.BookID, loan.BorrowerName, loan.DueDate, loan.Status)
	if err != nil {
		return err
	}

	loan.ID = int(id)
	loan.LoanDate = time.Now()
	loan.Status = "active"
	return nil
}

func (s *LoanService) GetLoans() ([]Loan, error) {
	query := "SELECT id, book_id, borrower_name, loan_date, due_date, return_date, status FROM loans"
	rows, err := s.db.Query(s.db.Rebind(query))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []Loan

	for rows.Next() {
		var l Loan
		var returnDate sql.NullTime
		err := rows.Scan(&l.ID, &l.BookID, &l.BorrowerName, &l.LoanDate, &l.DueDate, &returnDate, &l.Status)

		if err != nil {
			return nil, err
		}

		if returnDate.Valid {
			l.ReturnDate = returnDate.Time
		}
		loans = append(loans, l)
	}
	return loans, nil
}

func (s *LoanService) GetLoanByID(id int) (*Loan, error) {
	query := "SELECT id, book_id, borrower_name, loan_date, due_date, return_date, status FROM loans WHERE id = ?"
	row := s.db.QueryRow(s.db.Rebind(query), id)
	var l Loan
	var returnDate sql.NullTime
	err := row.Scan(&l.ID, &l.BookID, &l.BorrowerName, &l.LoanDate, &l.DueDate, &returnDate, &l.Status)

	if err != nil {
		return nil, err
	}

	if returnDate.Valid {
		l.ReturnDate = returnDate.Time
	}
	return &l, nil
}

func (s *LoanService) ReturnLoan(id int) error {
	query := "UPDATE loans SET status = 'returned', return_date = ? WHERE id = ? AND status = 'active'"
	result, err := s.db.Exec(s.db.Rebind(query), time.Now(), id)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *LoanService) GetOverdueLoans() ([]Loan, error) {
	query := "SELECT id, book_id, borrower_name, loan_date, due_date, return_date, status FROM loans WHERE due_date < ? AND status = 'active'"
	rows, err := s.db.Query(s.db.Rebind(query), time.Now())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []Loan

	for rows.Next() {
		var l Loan
		var returnDate sql.NullTime
		err := rows.Scan(&l.ID, &l.BookID, &l.BorrowerName, &l.LoanDate, &l.DueDate, &returnDate, &l.Status)
		if err != nil {
			return nil, err
		}
		if returnDate.Valid {
			l.ReturnDate = returnDate.Time
		}
		loans = append(loans, l)
	}
	return loans, nil
}
