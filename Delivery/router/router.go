package routers

import (
	"task_manager/Delivery/controllers"
	"task_manager/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctl *controllers.Controller, auth *infrastructure.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	// public
	r.POST("/register", ctl.Register)
	r.POST("/login", ctl.Login)

	// authenticated
	authGroup := r.Group("/")
	authGroup.Use(auth.AuthRequired())
	{
		authGroup.GET("/tasks", ctl.GetTasks)
		authGroup.GET("/tasks/:id", ctl.GetTaskByID)
	}

	// admin
	admin := r.Group("/")
	admin.Use(auth.AuthRequired(), auth.RequireAdmin())
	{
		admin.POST("/tasks", ctl.CreateTask)
		admin.PUT("/tasks/:id", ctl.UpdateTask)
		admin.DELETE("/tasks/:id", ctl.DeleteTask)
		admin.POST("/promote/:username", ctl.Promote)
	}
	return r
}
