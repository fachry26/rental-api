package routes

import (
	"github.com/gin-gonic/gin"
	"rental-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	maintenance := r.Group("/maintenance")
	{
		maintenance.POST("/", controllers.LogMaintenance)
		maintenance.GET("/:id", controllers.GetMaintenance)
		maintenance.GET("/", controllers.ListMaintenance)
	}

	rentals := r.Group("/rentals")
	{
		rentals.POST("/", controllers.CreateRental)
		rentals.GET("/:id", controllers.GetRental)
		rentals.GET("/", controllers.ListRentals)
		rentals.PUT("/:id/return", controllers.ReturnRental)
	}

	reviews := r.Group("/reviews")
	{
		reviews.POST("/", controllers.SubmitReview)
		reviews.GET("/:id", controllers.GetReview)
		reviews.GET("/", controllers.ListReviews)
		reviews.DELETE("/:id", controllers.DeleteReview)
	}

	users := r.Group("/users")
	{
		users.POST("/register", controllers.RegisterUser)
		users.POST("/login", controllers.LoginUser)
		users.GET("/:id", controllers.GetUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}
}
