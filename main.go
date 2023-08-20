package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS notes 
		(token TEXT, id INTEGER PRIMARY KEY, date TEXT, title TEXT, body TEXT)`)
	statement.Exec()
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	notes := []Note{}

	rows, err := db.Query("SELECT id, date, title, body FROM notes WHERE token = ?", token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Date, &note.Title, &note.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func main() {
	initDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/pull", getNotes).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
