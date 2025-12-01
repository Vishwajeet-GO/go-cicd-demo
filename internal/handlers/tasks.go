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

// CreateTask - Naya task banata hai
// POST /api/tasks
func CreateTask(c *gin.Context) {
	var task models.Task

	// Request body se data nikalo aur Task struct mein daalo
	if err := c.ShouldBindJSON(&task); err != nil {
		// Agar JSON galat format mein hai
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validation - Title khali nahi hona chahiye
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	// Database mein insert karo
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
	).Scan(&task.ID) // ID wapas le lo jo database ne generate kiya

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	// Timestamps set karo
	task.CreatedAt = now
	task.UpdatedAt = now

	// Success response bhejo
	c.JSON(http.StatusCreated, task)
}

// GetTasks - Saare tasks return karta hai
// GET /api/tasks
func GetTasks(c *gin.Context) {
	query := `
        SELECT id, title, description, completed, created_at, updated_at
        FROM tasks
        ORDER BY created_at DESC
    `

	// Database se query run karo
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tasks",
		})
		return
	}
	defer rows.Close() // Function khatam hone par rows close karo (cleanup)

	// Empty slice banao tasks store karne ke liye
	tasks := []models.Task{}

	// Har row ko iterate karo
	for rows.Next() {
		var task models.Task

		// Row se data nikalo aur struct mein daalo
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			continue // Agar error aaye to skip karo ye row
		}

		tasks = append(tasks, task) // Slice mein add karo
	}

	// JSON response bhejo
	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTask - Ek specific task return karta hai (ID se)
// GET /api/tasks/:id
func GetTask(c *gin.Context) {
	// URL se ID nikalo (string mein aata hai)
	idStr := c.Param("id")

	// String ko integer mein convert karo
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

	// Single row fetch karo
	err = database.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Task nahi mila
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

// UpdateTask - Task ko update karta hai
// PUT /api/tasks/:id
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

	// Update query
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

	// Check karo koi row update hui ya nahi
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

// DeleteTask - Task ko delete karta hai
// DELETE /api/tasks/:id
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

// HealthCheck - Server aur database health check
// GET /health
func HealthCheck(c *gin.Context) {
	// Database ping karo
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
