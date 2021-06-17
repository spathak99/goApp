package db

import(
	"database/sql"
)


var DB *sql.DB

//Initialize DB
func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}
