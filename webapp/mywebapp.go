package webapp

import "github.com/gin-gonic/gin"

func StartServer() {
	r := gin.Default()

	r.GET("/todos", listTodos)
	r.GET("/todos/:id", getTodo)
	r.POST("/todos", createTodo)
	r.POST("/todos/:id/done", markDone)

	r.Run(":8080")
}
