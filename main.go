package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Note struct {
	Note string `json:"note"`
	Id   string `json:"id"`
}

func main() {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	var notes = []*Note{}

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
			notes = append(notes, &newNote)
			w.WriteHeader(http.StatusCreated)
			w.Write(newNoteJSON)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		// get note by id
		id := r.URL.Path[len("/notes/"):]
		var note *Note
		for _, n := range notes {
			if n.Id == id {
				note = n
				break
			}
		}
		// check if note is empty
		if note == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("note not found"))
			return
		}
		switch r.Method {
		case "PUT":
			var updatedNote = r.FormValue("note")
			// check if note is empty
			if note.Note == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("note is empty"))
				return
			}
			// update note
			note.Note = updatedNote
			// convert updated note to json
			updatedNoteJSON, err := json.Marshal(note)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(updatedNoteJSON)
		case "DELETE":
			// delete note
			for i, n := range notes {
				if n.Id == id {
					notes = append(notes[:i], notes[i+1:]...)
					break
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("note deleted"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+port, nil)
}
