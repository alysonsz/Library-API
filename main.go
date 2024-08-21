package main

import (
	"fmt"
	"project_go/internal/services"
)

func main() {
	bookInformation := services.Book{
		ID:              12345,
		Title:           "Book Title",
		Author:          "Author Name",
		Pages:           200,
		Genre:           "Genre",
		PublicationYear: 1990,
	}
	fmt.Println(bookInformation.GetTitleAndAuthor())
}
