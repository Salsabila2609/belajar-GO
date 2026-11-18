package main

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

// Todo adalah bentuk data satu todo.
// Tag `json:"..."` menentukan nama field saat jadi JSON (untuk React).
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// ErrNotFound dikembalikan kalau todo dengan id tertentu tidak ada.
var ErrNotFound = errors.New("todo tidak ditemukan")

// openDB membuka (atau membuat) file database SQLite dan menyiapkan tabelnya.
func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite hanya aman satu penulis pada satu waktu

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT    NOT NULL,
		done INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// listTodos mengambil semua todo (READ semua).
func listTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT id, text, done FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{} // bukan nil, supaya JSON-nya [] dan bukan null
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

// createTodo menambah todo baru (CREATE).
func createTodo(db *sql.DB, text string) (Todo, error) {
	res, err := db.Exec("INSERT INTO todos (text) VALUES (?)", text)
	if err != nil {
		return Todo{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Todo{}, err
	}
	return Todo{ID: int(id), Text: text, Done: false}, nil
}

// updateTodo mengubah text dan done sebuah todo (UPDATE).
func updateTodo(db *sql.DB, id int, text string, done bool) (Todo, error) {
	res, err := db.Exec("UPDATE todos SET text = ?, done = ? WHERE id = ?", text, done, id)
	if err != nil {
		return Todo{}, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return Todo{}, err
	}
	if n == 0 {
		return Todo{}, ErrNotFound
	}
	return Todo{ID: id, Text: text, Done: done}, nil
}

// getTodo mengambil satu todo berdasarkan id.
func getTodo(db *sql.DB, id int) (Todo, error) {
	var t Todo
	err := db.QueryRow("SELECT id, text, done FROM todos WHERE id = ?", id).
		Scan(&t.ID, &t.Text, &t.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	return t, err
}

// deleteTodo menghapus todo (DELETE).
func deleteTodo(db *sql.DB, id int) error {
	res, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
