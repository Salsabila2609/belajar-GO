// Semua komunikasi dengan server Go ada di file ini.

export type Todo = {
  id: number
  text: string
  done: boolean
}

// Helper: kirim request, lempar Error kalau server membalas status gagal.
async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options)

  if (!res.ok) {
    // Server Go membalas error dalam bentuk {"error": "pesan"}
    let message = `Request gagal (${res.status})`
    try {
      const body = await res.json()
      if (body.error) message = body.error
    } catch {
      // body bukan JSON, pakai pesan bawaan
    }
    throw new Error(message)
  }

  // 204 No Content (hasil DELETE) tidak punya body
  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}

const jsonHeaders = { 'Content-Type': 'application/json' }

export function listTodos(): Promise<Todo[]> {
  return request<Todo[]>('/api/todos')
}

export function createTodo(text: string): Promise<Todo> {
  return request<Todo>('/api/todos', {
    method: 'POST',
    headers: jsonHeaders,
    body: JSON.stringify({ text }),
  })
}

export function updateTodo(
  id: number,
  changes: { text?: string; done?: boolean },
): Promise<Todo> {
  return request<Todo>(`/api/todos/${id}`, {
    method: 'PATCH',
    headers: jsonHeaders,
    body: JSON.stringify(changes),
  })
}

export function deleteTodo(id: number): Promise<void> {
  return request<void>(`/api/todos/${id}`, { method: 'DELETE' })
}
