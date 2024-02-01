package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	Id     string `json:"id" form:"id,omitempty" sql:"-"`
	Title  string `json:"title" form:"title" sql:"title"`
	Author string `json:"author" form:"author" sql:"author"`
}


var db *sql.DB


func main() {
	var err error
	db, err = sql.Open("sqlite3", "books.db")
	if err != nil{
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil{
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/create", getCreateBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", postBooks)
	router.GET("/", getHomePage)

	router.Run("localhost:8080")
}

func getCreateBooks(c *gin.Context){
	component := BooksCreateFormTemplate()
	component.Render(c, c.Writer)
}

func getHomePage(c *gin.Context){
	component := RootTemplate()
	component.Render(c, c.Writer)
}


func postBooks(c *gin.Context){
	var newBook Book
	if c.ShouldBind(&newBook) != nil{
		c.AbortWithError(http.StatusBadRequest, errors.New("the form is invalid"))
		return
	}
	stmt, err := db.Prepare("INSERT INTO books (author, title) VALUES (?, ?)")
	if err != nil{
		log.Fatal(err)
	}
	_, err = stmt.Exec(newBook.Author, newBook.Title)
	if err != nil{
		log.Fatal(err)
	}
	getBooks(c)
}

func getBook(c *gin.Context){
	id := c.Param("id")
	var book Book
	err := db.QueryRow("SELECT ROWID, title, author FROM books where rowid = ?", id).Scan(&book.Id, &book.Title, &book.Author)
	if err != nil{
		log.Fatal(err)
	}
	component := BooksSingleTemplate(&book)
	component.Render(c, c.Writer)
}

func getBooks(c *gin.Context) {
	books := retrieveBooksFromDatabase(c)

	component := BooksTemplate(books)
	component.Render(c, c.Writer)
}

func retrieveBooksFromDatabase(c *gin.Context) []Book {
	rows, err := db.Query("SELECT ROWID, title, author FROM books ORDER BY ROWID DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var newBook Book
		err := rows.Scan(&newBook.Id, &newBook.Title, &newBook.Author)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, newBook)
	}
	return books
}