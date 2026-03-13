package webapp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// in-memory store (replace with DB in real projects)
var todos = []Todo{
	{ID: 1, Text: "Learn Go basics", Done: true},
	{ID: 2, Text: "Learn Gin framework", Done: false},
}

// GET /todos
func listTodos(c *gin.Context) {
	resp := TodoListResponse{}
	for _, t := range todos {
		resp.Data = append(resp.Data, toTodoResponse(t))
	}
	c.JSON(http.StatusOK, resp)
}

// GET /todos/:id
func getTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, t := range todos {
		if t.ID == id {
			c.JSON(http.StatusOK, TodoItemResponse{Data: toTodoResponse(t)})
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
}

// POST /todos
func createTodo(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	newTodo := Todo{ID: len(todos) + 1, Text: req.Text}
	todos = append(todos, newTodo)
	c.JSON(http.StatusCreated, TodoItemResponse{Data: toTodoResponse(newTodo)})
}

// POST /todos/:id/done
func markDone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, t := range todos {
		if t.ID == id {
			todos[i].Done = true
			c.JSON(http.StatusOK, TodoItemResponse{Data: toTodoResponse(todos[i])})
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
}
