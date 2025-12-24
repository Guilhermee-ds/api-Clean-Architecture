package main

import (
	"api-clean-architecture/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func connectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:2717")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("not connected")
		log.Fatal(err)
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection = client.Database("taskbds").Collection("taskss")
	log.Println("Connected to MongoDB")
}

func createTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.CreatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func listTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)

	var tasks []models.Task

	if err = cursor.All(ctx, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func main() {
	connectDB()
	router := gin.Default()

	router.POST("/task", createTask)
	router.GET("/tasks", listTasks)

	router.Run(":8080")
}
