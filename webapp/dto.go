package webapp

// ---------- Requests ----------

type CreateTodoRequest struct {
	Text string `json:"text" binding:"required"`
}

// ---------- Responses ----------

type TodoResponse struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoListResponse struct {
	Data []TodoResponse `json:"data"`
}

type TodoItemResponse struct {
	Data TodoResponse `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// ---------- Mappers ----------

func toTodoResponse(t Todo) TodoResponse {
	return TodoResponse{ID: t.ID, Text: t.Text, Done: t.Done}
}
