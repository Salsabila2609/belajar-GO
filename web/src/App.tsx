import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { createTodo, deleteTodo, listTodos, updateTodo } from './api'
import type { Todo } from './api'
import TodoItem from './TodoItem'
import './App.css'

export default function App() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [text, setText] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // READ: ambil semua todo sekali saat halaman pertama kali tampil
  useEffect(() => {
    listTodos()
      .then(setTodos)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  // CREATE
  async function handleAdd(e: FormEvent) {
    e.preventDefault()
    const value = text.trim()
    if (!value) return
    try {
      const created = await createTodo(value)
      setTodos((prev) => [...prev, created])
      setText('')
      setError(null)
    } catch (err) {
      setError((err as Error).message)
    }
  }

  // UPDATE: centang / hapus centang
  async function handleToggle(todo: Todo) {
    try {
      const updated = await updateTodo(todo.id, { done: !todo.done })
      setTodos((prev) => prev.map((t) => (t.id === updated.id ? updated : t)))
      setError(null)
    } catch (err) {
      setError((err as Error).message)
    }
  }

  // UPDATE: ganti teks
  async function handleRename(todo: Todo, newText: string) {
    try {
      const updated = await updateTodo(todo.id, { text: newText })
      setTodos((prev) => prev.map((t) => (t.id === updated.id ? updated : t)))
      setError(null)
    } catch (err) {
      setError((err as Error).message)
    }
  }

  // DELETE
  async function handleDelete(todo: Todo) {
    try {
      await deleteTodo(todo.id)
      setTodos((prev) => prev.filter((t) => t.id !== todo.id))
      setError(null)
    } catch (err) {
      setError((err as Error).message)
    }
  }

  const remaining = todos.filter((t) => !t.done).length

  return (
    <main className="app">
      <h1>Todo</h1>

      <form className="add-form" onSubmit={handleAdd}>
        <input
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Apa yang mau dikerjakan?"
          aria-label="Todo baru"
        />
        <button type="submit" disabled={!text.trim()}>
          Tambah
        </button>
      </form>

      {error && (
        <p className="error" role="alert">
          {error}
        </p>
      )}

      {loading ? (
        <p className="muted">Memuat...</p>
      ) : todos.length === 0 ? (
        <p className="muted">Belum ada todo. Tambahkan satu di atas.</p>
      ) : (
        <>
          <ul className="list">
            {todos.map((todo) => (
              <TodoItem
                key={todo.id}
                todo={todo}
                onToggle={handleToggle}
                onRename={handleRename}
                onDelete={handleDelete}
              />
            ))}
          </ul>
          <p className="muted footer">
            {remaining} dari {todos.length} belum selesai
          </p>
        </>
      )}
    </main>
  )
}