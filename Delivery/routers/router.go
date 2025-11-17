package routers

import (
	"task_manager1/Delivery/controllers"
	"task_manager1/Infrastructure/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctl *controllers.Controller, authMw *auth.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", ctl.Register)
	r.POST("/login", ctl.Login)

	// Authenticated routes
	authGroup := r.Group("/")
	authGroup.Use(authMw.Handle())
	{
		authGroup.GET("/tasks", ctl.GetTasks)
		authGroup.GET("/tasks/:id", ctl.GetTaskByID)
	}

	// Admin routes
	admin := r.Group("/")
	admin.Use(authMw.Handle(), authMw.RequireAdmin())
	{
		admin.POST("/tasks", ctl.CreateTask)
		admin.PUT("/tasks/:id", ctl.UpdateTask)
		admin.DELETE("/tasks/:id", ctl.DeleteTask)
		admin.POST("/promote/:username", ctl.Promote)
	}

	return r
}
