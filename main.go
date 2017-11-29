package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "golang.org/x/tools/cmd/getgo/server"
)

type Bookmarks struct {
	Title       string
	Address     string
	Description string
	Time        time.Time
}

type SendBookmarksData struct {
	NewSlice   []Bookmarks
	TotalPages []int
}

func index(w http.ResponseWriter, r *http.Request) {
	var bookmarkSlice []Bookmarks
	// fmt.Fprintln(w, "Welcome to Bookmarks")
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	t := time.Now()
	timeZoneIdentifierForIndia := "Asia/Kolkata"
	india, _ := time.LoadLocation(timeZoneIdentifierForIndia)

	bookmark := Bookmarks{
		Title:       r.FormValue("title"),
		Address:     r.FormValue("address"),
		Description: r.FormValue("description"),
		Time:        t.In(india),
	}

	emailCookie, cookieError := r.Cookie("email")

	if cookieError == nil {
		rows, err := loadBookmarks(emailCookie.Value)
		if err != nil {
			panic(err)
		}
		var bookmarkTemp Bookmarks
		var email string
		for rows.Next() {
			err := rows.Scan(&bookmarkTemp.Title, &bookmarkTemp.Address, &bookmarkTemp.Description, &email)
			if err != nil {
				panic(err)
			} else {
				bookmarkSlice = append(bookmarkSlice, bookmarkTemp)
			}
		}
	}

	// var index int = int(page)
	// newbookmarkSlice := bookmarkSlice[10*(index-1) : (index*10)-1]

	if r.Method != http.MethodPost {
		var newBookmarkSlice []Bookmarks
		r.ParseForm()
		if bookmarkSlice == nil {
			tmpl.Execute(w, nil)
		} else {
			page := r.FormValue("page")
			if page != "" {
				//fmt.Println(page)
				index, err := strconv.ParseInt(page, 10, 64)
				//fmt.Println(index)
				if err != nil {
					fmt.Println(err)
				}
				if (index * 5) < int64(len(bookmarkSlice)) {
					newBookmarkSlice = bookmarkSlice[(index-1)*5 : (index * 5)]
				} else {
					newBookmarkSlice = bookmarkSlice[(index-1)*5:]
				}
			}
			//tmpl.Execute(w, newBookmarkSlice)
			sendData := SendBookmarksData{
				NewSlice:   newBookmarkSlice,
				TotalPages: totalPages(len(bookmarkSlice)/5 + 1),
			}
			fmt.Println(sendData.NewSlice)
			tmpl.Execute(w, sendData)
		}
		return
	} else {
		if bookmark.Address != "" {
			saveBookmark(bookmark, emailCookie.Value)
			bookmarkSlice = append(bookmarkSlice, bookmark)
		}
		var newBookmarkSlice []Bookmarks
		page := r.FormValue("page")
		if page != "" {
			//fmt.Println(page)
			index, err := strconv.ParseInt(page, 10, 64)
			//fmt.Println(index)
			if err != nil {
				fmt.Println(err)
			}

			if (index * 5) < int64(len(bookmarkSlice)) {
				newBookmarkSlice = bookmarkSlice[(index-1)*5 : (index*5)-1]
			} else {
				newBookmarkSlice = bookmarkSlice[(index-1)*5:]
			}
		}
		//fmt.Println(newBookmarkSlice)
		//tmpl.Execute(w, newBookmarkSlice)
		sendData := SendBookmarksData{
			NewSlice:   newBookmarkSlice,
			TotalPages: totalPages(len(bookmarkSlice)/5 + 1),
		}
		fmt.Println(sendData.NewSlice)
		tmpl.Execute(w, sendData)
	}

	// t := time.Now()
	// timeZoneIdentifierForIndia := "Asia/Kolkata"
	// india, _ := time.LoadLocation(timeZoneIdentifierForIndia)
	// fmt.Println(t.UTC())
	// fmt.Println(t.In(india))

	// fmt.Println(bookmarkSlice)
	//tmpl.Execute(w, bookmarkSlice)
}

func login(w http.ResponseWriter, r *http.Request) {
	confirmEmailCookie, _ := r.Cookie("email")
	confirmPasswordCookie, _ := r.Cookie("password")

	if confirmEmailCookie != nil && confirmPasswordCookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
	}
	r.ParseForm()

	user := User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if user.Password == selectAccount(user.Email) && user.Password != "" {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		emailCookie := http.Cookie{
			Name:    "email",
			Value:   user.Email,
			Expires: expiration,
		}
		passwordCookie := http.Cookie{
			Name:    "password",
			Value:   user.Password,
			Expires: expiration,
		}
		http.SetCookie(w, &emailCookie)
		http.SetCookie(w, &passwordCookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func totalPages(n int) (num []int) {
	var i int
	num = make([]int, n)
	for i = 0; i < n; i++ {
		num[i] = i + 1
	}
	return
}

func signup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/signup.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
	}

	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm_password")

	if email != "" && password != "" && confirm_password != "" {
		if password == confirm_password {
			createAccount(email, password)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	confirmEmailCookie, _ := r.Cookie("email")
	confirmPasswordCookie, _ := r.Cookie("password")
	emailCookie := http.Cookie{
		Name:    "email",
		Value:   confirmEmailCookie.Value,
		Expires: time.Now(),
	}
	passwordCookie := http.Cookie{
		Name:    "password",
		Value:   confirmPasswordCookie.Value,
		Expires: time.Now(),
	}

	http.SetCookie(w, &emailCookie)
	http.SetCookie(w, &passwordCookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func main() {
	// server := http.Server{
	// 	Addr: ":8080",
	// }

	// server := http.Server{
	// 	Addr: ":" + os.Getenv("PORT"),
	// }

	initDatabase()

	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static", files))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/logout", logout)

	// http.ListenAndServe(":3001", nil)
	server := &http.Server{
		Addr:    ":3001",
		Handler: mux,
	}
	server.ListenAndServe()
}
