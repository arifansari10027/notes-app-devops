package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Note struct {
	ID      int
	Content string
}

var (
	notes []Note
	id    int
	mu    sync.Mutex
)

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/delete", handleDelete)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server running on http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	mu.Lock()
	defer mu.Unlock()
	tmpl.Execute(w, notes)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content := r.FormValue("content")
		mu.Lock()
		id++
		notes = append(notes, Note{ID: id, Content: content})
		mu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idStr := r.FormValue("id")
		noteID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		mu.Lock()
		for i, note := range notes {
			if note.ID == noteID {
				notes = append(notes[:i], notes[i+1:]...)
				break
			}
		}
		mu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
