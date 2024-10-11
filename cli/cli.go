package cli

import (
	"fmt"
	"os"
	"project-go/internal/services"
	"strconv"
	"time"
)

type BookCLI struct {
	Service *services.BookService
}

func NewBookCLI(Service *services.BookService) *BookCLI {

	return &BookCLI{Service: Service}

}

func (cli *BookCLI) Run() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: books <command> [arguments]")
		return
	}

	command := os.Args[1]
	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: books search <book title>")
		}
		bookName := os.Args[2]
		cli.SearchBook(bookName)
	case "simulate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: books simulate <book_id> <book_id> <book_id> ...")
		}
		bookIDs := os.Args[2:]
		cli.SimulateReading(bookIDs)
	}
}

func (cli *BookCLI) SearchBook(nameBook string) {

	books, error := cli.Service.SearchBooksByName(nameBook)
	if error != nil {
		fmt.Println("Error searching books: ", error)
		return
	}
	if len(books) == 0 {
		fmt.Println("No books found")
		return
	}
	fmt.Printf("%d Books found\n", len(books))
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Genre: %s, Pages: %d, Publication Year: %d\n",
			book.ID, book.Title, book.Author, book.Genre, book.Pages, book.PublicationYear)
	}

}

func (cli *BookCLI) SimulateReading(BookIDs []string) {

	var bookIDs []int
	for _, idString := range BookIDs {
		id, error := strconv.Atoi(idString)
		if error != nil {
			fmt.Println("Invalid book ID: ", error)
			continue
		}
		bookIDs = append(bookIDs, id)
	}
	responses := cli.Service.SimulateMultipleReading(bookIDs, 2*time.Second)
	for _, response := range responses {
		fmt.Println(response)
	}

}
