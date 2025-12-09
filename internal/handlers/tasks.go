package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Vishwajeet-GO/go-cicd-demo/internal/database"
	"github.com/Vishwajeet-GO/go-cicd-demo/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	query := `
        INSERT INTO tasks (title, description, completed, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	now := time.Now()
	err := database.DB.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Completed,
		now,
		now,
	).Scan(&task.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	task.CreatedAt = now
	task.UpdatedAt = now

	c.JSON(http.StatusCreated, task)
}

func GetTasks(c *gin.Context) {
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM tasks
        ORDER BY created_at DESC
    `

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tasks",
		})
		return
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			continue
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

func GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	var task models.Task
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `

	err = database.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	query := `
        UPDATE tasks
        SET title = $1, description = $2, completed = $3, updated_at = $4
        WHERE id = $5
    `

	result, err := database.DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Completed,
		time.Now(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	task.ID = uint(id)
	task.UpdatedAt = time.Now()
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	query := `DELETE FROM tasks WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

func HealthCheck(c *gin.Context) {
	if err := database.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}
