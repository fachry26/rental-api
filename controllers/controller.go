package controllers

import (
	"net/http"
	"rental-api/models"
	"rental-api/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMaintenance(c *gin.Context) {
	var maintenance models.Maintenance
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	if err := models.CreateMaintenance(&maintenance); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to log maintenance: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusCreated, maintenance)
}

func GetMaintenance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid ID: "+err.Error())
		return
	}

	maintenance, err := models.GetMaintenanceByID(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Maintenance record not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, maintenance)
}

func ListMaintenance(c *gin.Context) {
	records, err := models.GetAllMaintenance()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch records: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, records)
}

func CreateRental(c *gin.Context) {
	var rental models.RentalHistory
	if err := c.ShouldBindJSON(&rental); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid rental data: "+err.Error())
		return
	}

	if rental.ReturnDate != nil && rental.RentalDate.After(*rental.ReturnDate) {
		utils.RespondError(c, http.StatusBadRequest, "Return date cannot be before rental date")
		return
	}

	if err := models.CreateRental(&rental); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create rental: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusCreated, rental)
}

func GetRental(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid rental ID: "+err.Error())
		return
	}

	rental, err := models.GetRentalByID(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Rental not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, rental)
}

func ListRentals(c *gin.Context) {
	rentals, err := models.GetAllRentals()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch rentals: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, rentals)
}

func ReturnRental(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid rental ID: "+err.Error())
		return
	}

	rental, err := models.GetRentalByID(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Rental not found")
		return
	}

	currentTime := time.Now()
	rental.ReturnDate = &currentTime

	if err := models.MarkAsReturned(id, rental.ReturnDate); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to return rental: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Rental returned successfully", "rental": rental})
}

func SubmitReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid review data: "+err.Error())
		return
	}

	if err := models.CreateReview(&review); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to submit review: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusCreated, review)
}

func GetReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid review ID: "+err.Error())
		return
	}

	review, err := models.GetReviewByID(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Review not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, review)
}

func ListReviews(c *gin.Context) {
	reviews, err := models.GetAllReviews()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch reviews: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, reviews)
}

func DeleteReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid review ID: "+err.Error())
		return
	}

	if err := models.DeleteReview(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to delete review: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data: "+err.Error())
		return
	}

	exists, err := models.CheckUserExists(user.Email)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to validate user: "+err.Error())
		return
	}
	if exists {
		utils.RespondError(c, http.StatusConflict, "User already exists")
		return
	}

	if err := models.CreateUser(&user); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data: "+err.Error())
		return
	}

	user, err := models.AuthenticateUser(loginData.Email, loginData.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to generate token: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid user ID: "+err.Error())
		return
	}

	user, err := models.GetUserByID(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := models.UpdateUser(id, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := models.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func CreateMachine(c *gin.Context) {
	var machine models.MesinBor
	if err := c.ShouldBindJSON(&machine); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid machine data: "+err.Error())
		return
	}

	if err := models.CreateMachine(&machine); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create machine: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusCreated, machine)
}

func GetMachine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid machine ID: "+err.Error())
		return
	}

	machine, err := models.GetMachineByID(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "Machine not found")
		return
	}

	utils.RespondJSON(c, http.StatusOK, machine)
}

func ListMachines(c *gin.Context) {
	machines, err := models.GetMachines()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch machines: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, machines)
}

func UpdateMachine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid machine ID: "+err.Error())
		return
	}

	var updatedMachine models.MesinBor
	if err := c.ShouldBindJSON(&updatedMachine); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid machine data: "+err.Error())
		return
	}

	if err := models.UpdateMachine(id, &updatedMachine); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update machine: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Machine updated successfully"})
}

func DeleteMachine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid machine ID: "+err.Error())
		return
	}

	if err := models.DeleteMachine(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to delete machine: "+err.Error())
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{"message": "Machine deleted successfully"})
}
