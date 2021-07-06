package URL

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func shorting() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

type Result struct {
	Link   string
	Code   string
	Status string
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/index.html")
	result := Result{}
	if r.Method == "POST" {
		if !isValidUrl(r.FormValue("s")) {
			fmt.Println("Что-то не так")
			result.Status = "Ссылка имеет неправильный формат!"
			result.Link = ""
		} else {
			result.Link = r.FormValue("s")
			result.Code = shorting()
			db, err := sql.Open("sqlite3", "project.db")
			if err != nil {
				panic(err)
			}
			defer db.Close()
			db.Exec("insert into links (link, short) values ($1, $2)", result.Link, result.Code)
			result.Status = "Сокращение было выполнено успешно"
		}
	}
	templ.Execute(w, result)
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	var link string
	vars := mux.Vars(r)
	db, err := sql.Open("sqlite3", "project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows := db.QueryRow("select link from links where short=$1 limit 1", vars["key"])
	rows.Scan(&link)
	fmt.Fprintf(w, "<script>location='%s';</script>", link)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/to/{key}", redirectTo)
	log.Fatal(http.ListenAndServe(":8000", router))
}
