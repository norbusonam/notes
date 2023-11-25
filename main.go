package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Note struct {
	Note string `json:"note"`
	Id   string `json:"id"`
}

func main() {
	var notes = []Note{}

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// convert notes to json list
			notesJSON, err := json.Marshal(notes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(notesJSON)
		case "POST":
			note := r.FormValue("note")
			// check if note is empty
			if note == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("note is empty"))
				return
			}
			// create new note
			newNote := Note{
				Note: note,
				Id:   uuid.New().String(),
			}
			// convert new note to json
			newNoteJSON, err := json.Marshal(newNote)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			// persist new note by appending to notes list
			notes = append(notes, newNote)
			w.WriteHeader(http.StatusCreated)
			w.Write(newNoteJSON)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)

}
