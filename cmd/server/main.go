package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/hawkaii/website-checker.git/internal/handlers"
	"github.com/hawkaii/website-checker.git/internal/storage"
	"github.com/hawkaii/website-checker.git/pkg/checker"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize Gin router
	router := gin.Default()

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // fallback
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

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
	log.Printf("Using MongoDB URI: %s", mongoURI)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
