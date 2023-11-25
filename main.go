package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	var notes = []string{}

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// convert notes to
			notesJSON, err := json.Marshal(notes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(notesJSON)
		case "POST":
			// check for empty note
			var note = r.FormValue("note")
			if note == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("note is empty"))
				return
			}
			notes = append(notes, note)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(note))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)

}
