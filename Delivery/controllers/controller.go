package controllers

import (
	"context"
	"net/http"

	"task_manager1/Domain"
	"task_manager1/Infrastructure/auth"
	"task_manager1/Usecases"

	"github.com/gin-gonic/gin"
)

// Controller holds Usecases and services
type Controller struct {
	UserUC *Usecases.UserUsecase
	TaskUC *Usecases.TaskUsecase
	JWT    *auth.JWTService
}

// NewController constructs controller
func NewController(userUC *Usecases.UserUsecase, taskUC *Usecases.TaskUsecase, jwt *auth.JWTService) *Controller {
	return &Controller{UserUC: userUC, TaskUC: taskUC, JWT: jwt}
}

// Register endpoint
func (ctl *Controller) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}
	ctx := context.Background()
	u, err := ctl.UserUC.Register(ctx, body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// generate token
	token, err := ctl.JWT.GenerateToken(u.Username, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"username": u.Username, "role": u.Role, "token": token})
}

// Login endpoint
func (ctl *Controller) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}
	ctx := context.Background()
	u, err := ctl.UserUC.Authenticate(ctx, body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := ctl.JWT.GenerateToken(u.Username, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": u.Username, "role": u.Role, "token": token})
}

// Promote endpoint (admin)
func (ctl *Controller) Promote(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	ctx := context.Background()
	updated, err := ctl.UserUC.Promote(ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to promote"})
		return
	}
	if updated.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": updated.Username, "role": updated.Role})
}

// GetTasks (authenticated)
func (ctl *Controller) GetTasks(c *gin.Context) {
	ctx := context.Background()
	tasks, err := ctl.TaskUC.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}
	resp := []Domain.TaskResponse{}
	for _, t := range tasks {
		id := ""
		if !t.ID.IsZero() {
			id = t.ID.Hex()
		}
		resp = append(resp, Domain.TaskResponse{
			ID:          id,
			Title:       t.Title,
			Description: t.Description,
			DueDate:     t.DueDate,
			Status:      t.Status,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// GetTaskByID (authenticated)
func (ctl *Controller) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	t, err := ctl.TaskUC.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task"})
		return
	}
	if t.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, Domain.TaskResponse{
		ID:          t.ID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		Status:      t.Status,
	})
}

// CreateTask (admin)
func (ctl *Controller) CreateTask(c *gin.Context) {
	var input Domain.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	ctx := context.Background()
	created, err := ctl.TaskUC.Create(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, Domain.TaskResponse{
		ID:          created.ID.Hex(),
		Title:       created.Title,
		Description: created.Description,
		DueDate:     created.DueDate,
		Status:      created.Status,
	})
}

// UpdateTask (admin)
func (ctl *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input Domain.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	ctx := context.Background()
	updated, err := ctl.TaskUC.Update(ctx, id, input)
	if err != nil {
		if err.Error() == "no fields to update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	if updated.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, Domain.TaskResponse{
		ID:          updated.ID.Hex(),
		Title:       updated.Title,
		Description: updated.Description,
		DueDate:     updated.DueDate,
		Status:      updated.Status,
	})
}

// DeleteTask (admin)
func (ctl *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	ok, err := ctl.TaskUC.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
