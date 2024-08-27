package main

import (
	"log"
	"time"

	"github.com/dagota12/Loan-Tracker/api/route"
	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the app
	app := bootstrap.App()

	// Get the environment variables
	env := app.Env

	// Connect to the database
	db := app.Mongo.Database(env.DBName)

	// Close the database connection when the main function is done
	defer app.CloseDBConnection()

	// Set the timeout for the context of the request
	timeout := time.Duration(env.ContextTimeout) * time.Second
	log.Println("[main] context timeout", timeout)
	// Initialize the gin
	gin := gin.Default()

	// Setup the routes
	route.Setup(env, timeout, db, gin)

	// Run the server
	gin.Run(env.ServerAddress)

}
