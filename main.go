package main

import (
	"log"
	"rental-api/controllers"
	"rental-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	err := models.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	defer models.CloseDatabase()

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", controllers.RegisterUser)
		userGroup.POST("/login", controllers.LoginUser)
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}

	maintenanceGroup := r.Group("/maintenance")
	{
		maintenanceGroup.POST("/", controllers.LogMaintenance)
		maintenanceGroup.GET("/:id", controllers.GetMaintenance)
		maintenanceGroup.GET("/", controllers.ListMaintenance)
	}

	rentalGroup := r.Group("/rentals")
	{
		rentalGroup.POST("/", controllers.CreateRental)
		rentalGroup.GET("/:id", controllers.GetRental)
		rentalGroup.GET("/", controllers.ListRentals)
		rentalGroup.PUT("/:id/return", controllers.ReturnRental)
	}

	reviewGroup := r.Group("/reviews")
	{
		reviewGroup.POST("/", controllers.SubmitReview)
		reviewGroup.GET("/:id", controllers.GetReview)
		reviewGroup.GET("/", controllers.ListReviews)
		reviewGroup.DELETE("/:id", controllers.DeleteReview)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
