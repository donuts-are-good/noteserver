package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/syslog"
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

const servePort = ":4096"

var db *sql.DB
var logger *syslog.Writer

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
		logger.Err("Failed to open database: " + err.Error())
		log.Fatal(err)
	}

	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS notes 
		(token TEXT, id INTEGER PRIMARY KEY, date TEXT, title TEXT, body TEXT)`)
	statement.Exec()
}

func syncNotes(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if token == "" || len(token) != 64 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		notes := []Note{}

		rows, err := db.Query("SELECT id, date, title, body FROM notes WHERE token = ?", token)
		if err != nil {
			logger.Err("Database query failed: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var note Note
			if err := rows.Scan(&note.ID, &note.Date, &note.Title, &note.Body); err != nil {
				logger.Err("Failed to scan database rows: " + err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			notes = append(notes, note)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)

	case "POST":
		var notes []Note
		err := json.NewDecoder(r.Body).Decode(&notes)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO notes(token, date, title, body) VALUES (?, ?, ?, ?)")
		if err != nil {
			logger.Err("Failed to prepare statement: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for _, note := range notes {
			_, err := stmt.Exec(token, note.Date, note.Title, note.Body)
			if err != nil {
				logger.Err("Failed to execute statement: " + err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	var err error
	logger, err = syslog.New(syslog.LOG_ERR|syslog.LOG_LOCAL0, "notes-app")
	if err != nil {
		log.Fatal("Failed to initialize syslog: ", err)
	}
	defer logger.Close()

	initDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/sync", syncNotes).Methods("GET", "POST")

	http.Handle("/", r)

	log.Println("Alive! Serving on " + servePort)
	http.ListenAndServe(servePort, nil)
}
