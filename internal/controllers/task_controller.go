package controllers

import (
	"net/http"
	"working-day-api/internal/domain"
	"working-day-api/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	TaskService *services.TaskService
}

func NewTaskController(taskService *services.TaskService) *TaskController {
	return &TaskController{TaskService: taskService}
}

func (ctrl *TaskController) GetTaskHandler(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
		return
	}

	task, err := ctrl.TaskService.GetTask(id, userID.(uint), userRole.(string))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you do not have permission to view this task" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ctrl *TaskController) CreateTaskHandler(c *gin.Context) {
	var task domain.Task

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	err := ctrl.TaskService.CreateTask(&task, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"status":  http.StatusCreated,
	})
}

func (ctrl *TaskController) UpdateTaskHandler(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updatedTask domain.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	err := ctrl.TaskService.UpdateTask(id, userID.(uint), &updatedTask)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you do not have permission to update this task" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
	})
}

func (ctrl *TaskController) DeleteTaskHandler(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
		return
	}

	err := ctrl.TaskService.DeleteTask(id, userID.(uint), userRole.(string))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you do not have permission to delete this task" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

func (ctrl *TaskController) GetTasksHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
		return
	}

	tasks, err := ctrl.TaskService.GetTasks(userID.(uint), userRole.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No tasks found",
			"tasks":   []domain.Task{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks retrieved successfully",
		"tasks":   tasks,
	})
}
