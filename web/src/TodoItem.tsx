import { useState } from 'react'
import type { Todo } from './api'

type Props = {
  todo: Todo
  onToggle: (todo: Todo) => void
  onRename: (todo: Todo, text: string) => void
  onDelete: (todo: Todo) => void
}

export default function TodoItem({ todo, onToggle, onRename, onDelete }: Props) {
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(todo.text)

  function startEdit() {
    setDraft(todo.text)
    setEditing(true)
  }

  function save() {
    const text = draft.trim()
    setEditing(false)
    // Hanya kirim ke server kalau benar-benar berubah dan tidak kosong
    if (text && text !== todo.text) onRename(todo, text)
  }

  return (
    <li className={todo.done ? 'todo done' : 'todo'}>
      <input
        type="checkbox"
        checked={todo.done}
        onChange={() => onToggle(todo)}
        aria-label={`Tandai selesai: ${todo.text}`}
      />

      {editing ? (
        <input
          className="edit-input"
          value={draft}
          autoFocus
          onChange={(e) => setDraft(e.target.value)}
          onBlur={save}
          onKeyDown={(e) => {
            if (e.key === 'Enter') save()
            if (e.key === 'Escape') setEditing(false)
          }}
        />
      ) : (
        <span className="text" onDoubleClick={startEdit}>
          {todo.text}
        </span>
      )}

      <div className="actions">
        <button type="button" className="ghost" onClick={startEdit}>
          Edit
        </button>
        <button type="button" className="danger" onClick={() => onDelete(todo)}>
          Hapus
        </button>
      </div>
    </li>
  )
}
