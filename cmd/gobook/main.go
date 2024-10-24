package main

import (
	"database/sql"
	"net/http"
	"os"
	"project-go/internal/cli"
	"project-go/internal/services"
	"project-go/internal/web"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	database, error := sql.Open("sqlite3", "./books.db")
	if error != nil {

		panic(error)

	}
	defer database.Close()

	BookService := services.NewBookService(database)
	BookHandlers := web.NewBookHandlers(BookService)

	if len(os.Args) > 1 && (os.Args[1] == "simulate" || os.Args[1] == "search") {
		bookCLI := cli.NewBookCLI(BookService)
		bookCLI.Run()
		return
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /books", BookHandlers.GetBooks)
	router.HandleFunc("POST /books", BookHandlers.CreateBook)
	router.HandleFunc("GET /books/{id}", BookHandlers.GetBookByID)
	router.HandleFunc("PUT /books/{id}", BookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", BookHandlers.DeleteBook)

	http.ListenAndServe(":8080", router)

}
