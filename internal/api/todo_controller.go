package api

import (
	"errors"
	"strconv"

	"to-do-list/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type TodoController struct {
	svc *service.TodoService
}

func NewTodoController(svc *service.TodoService) *TodoController {
	return &TodoController{svc: svc}
}

func (t *TodoController) Register(r *gin.Engine) {
	r.GET("/todo", t.list)
	r.POST("/todo", t.create)
	r.PATCH("/todo/:id", t.update)
	r.DELETE("/todo/:id", t.delete)
}

func (t *TodoController) list(c *gin.Context) {
	out, err := t.svc.List(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, out)
}

func (t *TodoController) create(c *gin.Context) {
	var in struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid body"})
		return
	}
	out, err := t.svc.Create(c.Request.Context(), in.Title)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, out)
}

func (t *TodoController) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid id"})
		return
	}
	var in struct {
		Title *string `json:"title"`
		Done  *bool   `json:"done"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid body"})
		return
	}
	out, err := t.svc.Update(c.Request.Context(), id, in.Title, in.Done)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatusJSON(404, gin.H{"error": "not found"})
			return
		}
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, out)
}

func (t *TodoController) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid id"})
		return
	}
	ok, err := t.svc.Delete(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(404, gin.H{"error": "not found"})
		return
	}
	c.Status(204)
}
