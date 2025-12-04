package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHealthCheckEndpoint - Tests if health endpoint responds
func TestHealthCheckEndpoint(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Simple health check without DB dependency
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Test request
	req, _ := http.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assertions
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response["status"])
	}
}

// TestTaskCreationValidation - Tests input validation
func TestTaskCreationValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock handler for testing validation
	router.POST("/api/tasks", func(c *gin.Context) {
		var task struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Completed   bool   `json:"completed"`
		}

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate title
		if task.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
			return
		}

		c.JSON(http.StatusCreated, task)
	})

	// Test with empty title (should fail)
	task := map[string]interface{}{
		"title":       "",
		"description": "Test",
		"completed":   false,
	}

	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Should return 400 Bad Request
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty title, got %d", resp.Code)
	}
}

// TestTaskCreationSuccess - Tests successful task creation
func TestTaskCreationSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/tasks", func(c *gin.Context) {
		var task struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Completed   bool   `json:"completed"`
		}

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		if task.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
			return
		}

		// Mock successful creation
		task_response := map[string]interface{}{
			"id":          1,
			"title":       task.Title,
			"description": task.Description,
			"completed":   task.Completed,
		}

		c.JSON(http.StatusCreated, task_response)
	})

	// Test with valid data
	task := map[string]interface{}{
		"title":       "Test Task",
		"description": "Testing API",
		"completed":   false,
	}

	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Should return 201 Created
	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response["id"] != float64(1) {
		t.Errorf("Expected id 1, got %v", response["id"])
	}

	if response["title"] != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %v", response["title"])
	}
}

// TestJSONValidation - Tests JSON parsing
func TestJSONValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/tasks", func(c *gin.Context) {
		var task map[string]interface{}

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Valid JSON"})
	})

	// Test with invalid JSON
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", resp.Code)
	}
}
