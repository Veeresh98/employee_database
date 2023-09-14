package main

import (
	context2 "context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

var collect *mongo.Collection

type Employee struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstname"`
	SecondName string `json:"secondname"`
	Email      string `json:"email"`
}

func getEmployee(c *gin.Context) {

	option, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collect.Find(option, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error for creating the collection"})
		return
	}
	defer func(cursor *mongo.Cursor, ctx context2.Context) {
		err = cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, option)

	var employees []Employee

	for cursor.Next(option) {
		var employee Employee
		if err = cursor.Decode(&employee); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error in decoding the JSON file "})
			return
		}

		employees = append(employees, employee)
	}

	c.JSON(http.StatusOK, employees)

}

func getEmployeeByID(c *gin.Context) {

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var employee Employee

	err := collect.FindOne(ctx, Employee{ID: id}).Decode(&employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "employee id not found"})
		return
	}
	c.JSON(http.StatusOK, employee)

}

func createEmployee(c *gin.Context) {

	var employee Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to bind"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := collect.InsertOne(ctx, employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert it to teh collection"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func updateEmployee(c *gin.Context) {

	id := c.Param("id")
	var updateEmployee Employee

	if err := c.ShouldBindJSON(&updateEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error binding the JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collect.UpdateOne(ctx, Employee{ID: id}, bson.M{"$set": updateEmployee})
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"error": "updated the employee"})
	}

	c.JSON(http.StatusOK, updateEmployee)
}
func deleteEmployee(c *gin.Context) {

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := collect.DeleteOne(ctx, Employee{ID: id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting the id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "able to delete the id "})
}

func main() {

	ctx := context.Background()
	clientOpt := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOpt)

	if err != nil {
		return
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err = client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)

	router := gin.Default()

	router.GET("/employee", getEmployee)
	router.GET("/employee/:id", getEmployeeByID)
	router.POST("/employee", createEmployee)
	router.PUT("/employee/:id", updateEmployee)
	router.DELETE("/employee/:id", deleteEmployee)

	err = router.Run("localhost:8000")
	if err != nil {
		return
	}

	// Created the database and the collection.
	collect = client.Database("employeeDatabase").Collection("employee")

}
