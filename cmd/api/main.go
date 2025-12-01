package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vishwajeet-GO/go-cicd-demo/internal/database"
	"github.com/Vishwajeet-GO/go-cicd-demo/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env file load karo (local development ke liye)
	// Production mein ye file nahi hogi, environment variables directly set honge
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database se connect karo
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close() // Program exit hone par connection close karo

	// Gin router setup
	router := setupRouter()

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// HTTP server banao
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second, // Request padhne ka max time
		WriteTimeout: 10 * time.Second, // Response likhne ka max time
		IdleTimeout:  60 * time.Second, // Idle connection ka max time
	}

	// Server ko separate goroutine mein chalao
	// Goroutine = lightweight thread (background task)
	go func() {
		log.Printf("🚀 Server starting on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown setup
	// Jab Ctrl+C press karo ya system shutdown ho, to properly cleanup karo
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Wait karo signal ke liye (blocking operation)

	log.Println("🛑 Shutting down server...")

	// 5 second ka timeout do ongoing requests ko complete karne ke liye
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

// setupRouter - Saare routes define karta hai
func setupRouter() *gin.Engine {
	// Production mode mein kam logs (better performance)
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default() // Default middleware ke saath (logger, recovery)

	// CORS middleware (agar frontend alag hai to)
	router.Use(corsMiddleware())

	// Health check endpoint (CI/CD ke liye zaroori)
	router.GET("/health", handlers.HealthCheck)

	// API routes group
	api := router.Group("/api")
	{
		// Tasks endpoints
		api.POST("/tasks", handlers.CreateTask)       // Naya task
		api.GET("/tasks", handlers.GetTasks)          // Saare tasks
		api.GET("/tasks/:id", handlers.GetTask)       // Ek task
		api.PUT("/tasks/:id", handlers.UpdateTask)    // Task update
		api.DELETE("/tasks/:id", handlers.DeleteTask) // Task delete
	}

	return router
}

// corsMiddleware - Cross-Origin Resource Sharing
// Frontend aur backend alag domain par ho to ye allow karta hai
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Preflight request handle karo
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next() // Agle handler pe jao
	}
}
