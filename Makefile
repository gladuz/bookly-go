db: 
	sqlite3 -init database.sql books.db	 ""

run:
	templ generate
	go run .