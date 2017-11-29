package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	// "github.com/ziutek/mymysql/mysql"
	// _ "github.com/ziutek/mymysql/native"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

var Db *sql.DB

func initDatabase() {
	var err error
	// Db, err = sql.Open("mysql", "root@/bookmarks_db")
	Db, err = sql.Open("mysql", "usb4qllzicum66zz:wyleAL0rYiBw5vdNEE6@/bcxxphn9v")

	// Db := mysql.New("tcp",
	// 	"",
	// 	os.Getenv("MYSQL_ADDON_HOST") + ":" + os.Getenv("MYSQL_ADDON_PORT"),
	// 	os.Getenv("MYSQL_ADDON_USER"),
	// 	os.Getenv("MYSQL_ADDON_PASSWORD"),
	// 	os.Getenv("MYSQL_ADDON_DB")
	if err != nil {
		panic(err)
	}
}

func selectAccount(email string) (password string) {
	var err2 error
	var user User
	users, err2 := Db.Query("select Password from user_accounts where email = ?", email)

	if err2 != nil {
		panic(err2)
	}

	for users.Next() {
		result := users.Scan(&user.Password)
		fmt.Println(user.Password)
		password = user.Password
		if result != nil {
			panic(result)
		}
		return
	}
	return
}

func saveBookmark(bookmark Bookmarks, email string) {
	_, err := Db.Exec("insert into user_bookmarks set title=?, address=?, description=?, email=?", bookmark.Title,
		bookmark.Address, bookmark.Description, email)
	if err != nil {
		panic(err)
	}
}

func loadBookmarks(email string) (rows *sql.Rows, err error) {
	rows, err = Db.Query("Select title, address, description, email from user_bookmarks where email=?", email)
	return
}

func createAccount(email string, password string) {
	_, err := Db.Exec("insert into user_accounts set email=?, password=?", email, password)
	if err != nil {
		panic(err)
	}
}
