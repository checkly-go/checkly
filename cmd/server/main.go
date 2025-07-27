package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/hawkaii/website-checker.git/internal/handlers"
	"github.com/hawkaii/website-checker.git/internal/storage"
	"github.com/hawkaii/website-checker.git/pkg/checker"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Connect to MongoDB (adjust the URI as needed)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}
	db := client.Database("website_checker")

	// Create storage layer repositories.
	// User repository is set up for future use.
	userRepo := storage.NewUserRepository(db)
	checkRepo := storage.NewCheckRepository(db)

	// Create a new Checker instance from the core library package.
	// Ensure that the NewChecker() exists in pkg/checker.
	chk := checker.NewChecker()

	// Create the Service object that holds the checker and storage repositories.
	service := &handlers.Service{
		Checker:   chk,
		UserRepo:  userRepo,
		CheckRepo: checkRepo,
	}

	// Define API endpoints (all endpoints are currently unauthenticated).
	api := router.Group("/api/v1")
	{
		api.POST("/check", service.SubmitCheck)
		api.GET("/check/:id", service.GetCheck)
		api.GET("/check/:id/report", service.GetCheckReport)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	// Configure and start the HTTP server.
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting API server on port 8080...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
