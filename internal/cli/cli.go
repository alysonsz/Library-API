package cli

import (
	"fmt"
	"os"
	"project-go/internal/services"
	"strconv"
	"time"
)

type BookCLI struct {
	bookService   *services.BookService
	authorService *services.AuthorService
	loanService   *services.LoanService
}

func NewBookCLI(bookService *services.BookService, authorService *services.AuthorService, loanService *services.LoanService) *BookCLI {
	return &BookCLI{
		bookService:   bookService,
		authorService: authorService,
		loanService:   loanService,
	}
}

func (cli *BookCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gobook <command> [arguments]")
		fmt.Println("Commands: search, simulate, authors, loans")
		return
	}

	command := os.Args[1]
	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gobook search <book title>")
			return
		}
		cli.SearchBook(os.Args[2])
	case "simulate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gobook simulate <book_id> <book_id> ...")
			return
		}
		cli.SimulateReading(os.Args[2:])
	case "authors":
		cli.handleAuthors()
	case "loans":
		cli.handleLoans()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func (cli *BookCLI) SearchBook(nameBook string) {
	books, err := cli.bookService.SearchBooksByName(nameBook)
	if err != nil {
		fmt.Println("Error searching books: ", err)
		return
	}
	if len(books) == 0 {
		fmt.Println("No books found")
		return
	}
	fmt.Printf("%d Books found\n", len(books))
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Genre: %s, Pages: %d, Year: %d\n",
			book.ID, book.Title, book.Author, book.Genre, book.Pages, book.PublicationYear)
	}
}

func (cli *BookCLI) SimulateReading(BookIDs []string) {
	var ids []int
	for _, idString := range BookIDs {
		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println("Invalid book ID: ", err)
			continue
		}
		ids = append(ids, id)
	}
	responses := cli.bookService.SimulateMultipleReading(ids, 2*time.Second)
	for _, response := range responses {
		fmt.Println(response)
	}
}

func (cli *BookCLI) handleAuthors() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gobook authors <list|create|update|delete>")
		return
	}
	sub := os.Args[2]
	switch sub {
	case "list":
		authors, err := cli.authorService.GetAuthors()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, a := range authors {
			fmt.Printf("ID: %d, Name: %s, Bio: %s, BirthYear: %d\n", a.ID, a.Name, a.Bio, a.BirthYear)
		}
	case "create":
		if len(os.Args) < 4 {
			fmt.Println("Usage: gobook authors create <name> [bio] [birthYear]")
			return
		}

		author := services.Author{Name: os.Args[3]}
		if len(os.Args) > 4 {
			author.Bio = os.Args[4]
		}

		if len(os.Args) > 5 {
			author.BirthYear, _ = strconv.Atoi(os.Args[5])
		}

		if err := cli.authorService.CreateAuthor(&author); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Author created with ID: %d\n", author.ID)
	case "delete":
		if len(os.Args) < 4 {
			fmt.Println("Usage: gobook authors delete <id>")
			return
		}
		id, _ := strconv.Atoi(os.Args[3])

		if err := cli.authorService.DeleteAuthor(id); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Author deleted")
	}
}

func (cli *BookCLI) handleLoans() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: gobook loans <list|overdue|return>")
		return
	}

	sub := os.Args[2]
	switch sub {
	case "list":
		loans, err := cli.loanService.GetLoans()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, l := range loans {
			fmt.Printf("ID: %d, BookID: %d, Borrower: %s, Status: %s, Due: %s\n",
				l.ID, l.BookID, l.BorrowerName, l.Status, l.DueDate.Format(time.DateTime))
		}
	case "overdue":
		loans, err := cli.loanService.GetOverdueLoans()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("%d overdue loans\n", len(loans))

		for _, l := range loans {
			fmt.Printf("ID: %d, BookID: %d, Borrower: %s, Due: %s\n",
				l.ID, l.BookID, l.BorrowerName, l.DueDate.Format(time.DateTime))
		}
	case "return":
		if len(os.Args) < 4 {
			fmt.Println("Usage: gobook loans return <loan_id>")
			return
		}
		id, _ := strconv.Atoi(os.Args[3])

		if err := cli.loanService.ReturnLoan(id); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Loan returned")
	}
}
