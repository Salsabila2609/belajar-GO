package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// writeJSON mengirim balasan JSON dengan status code tertentu.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// parseID membaca {id} dari URL dan mengubahnya jadi angka.
func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func main() {
	db, err := openDB("todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/todos", handleList(db))
	mux.HandleFunc("POST /api/todos", handleCreate(db))
	mux.HandleFunc("PATCH /api/todos/{id}", handleUpdate(db))
	mux.HandleFunc("DELETE /api/todos/{id}", handleDelete(db))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := listTodos(db)
		if err != nil {
			log.Println("list:", err)
			writeError(w, http.StatusInternalServerError, "gagal membaca data")
			return
		}
		writeJSON(w, http.StatusOK, todos)
	}
}

func handleCreate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "JSON tidak valid")
			return
		}
		in.Text = strings.TrimSpace(in.Text)
		if in.Text == "" {
			writeError(w, http.StatusBadRequest, "text wajib diisi")
			return
		}

		todo, err := createTodo(db, in.Text)
		if err != nil {
			log.Println("create:", err)
			writeError(w, http.StatusInternalServerError, "gagal menyimpan")
			return
		}
		writeJSON(w, http.StatusCreated, todo)
	}
}

// handleUpdate memakai PATCH: field yang tidak dikirim tidak diubah.
// Pointer (*string, *bool) dipakai supaya bisa membedakan
// "tidak dikirim" (nil) dari "dikirim kosong/false".
func handleUpdate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id tidak valid")
			return
		}

		var in struct {
			Text *string `json:"text"`
			Done *bool   `json:"done"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "JSON tidak valid")
			return
		}

		current, err := getTodo(db, id)
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo tidak ditemukan")
			return
		}
		if err != nil {
			log.Println("get:", err)
			writeError(w, http.StatusInternalServerError, "gagal membaca data")
			return
		}

		if in.Text != nil {
			t := strings.TrimSpace(*in.Text)
			if t == "" {
				writeError(w, http.StatusBadRequest, "text tidak boleh kosong")
				return
			}
			current.Text = t
		}
		if in.Done != nil {
			current.Done = *in.Done
		}

		updated, err := updateTodo(db, id, current.Text, current.Done)
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo tidak ditemukan")
			return
		}
		if err != nil {
			log.Println("update:", err)
			writeError(w, http.StatusInternalServerError, "gagal mengubah")
			return
		}
		writeJSON(w, http.StatusOK, updated)
	}
}

func handleDelete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id tidak valid")
			return
		}

		err = deleteTodo(db, id)
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo tidak ditemukan")
			return
		}
		if err != nil {
			log.Println("delete:", err)
			writeError(w, http.StatusInternalServerError, "gagal menghapus")
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204: sukses, tanpa isi
	}
}
