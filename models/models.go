package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("rental.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting SQL DB instance:", err)
		return
	}
	sqlDB.Close()
}

type Maintenance struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MachineID uint      `json:"machine_id"`
	Issue     string    `json:"issue"`
	Fixed     bool      `json:"fixed"`
	FixedAt   time.Time `json:"fixed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RentalHistory struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	UserID     uint       `json:"user_id"`
	MachineID  uint       `json:"machine_id"`
	RentalDate time.Time  `json:"rental_date"`
	ReturnDate *time.Time `json:"return_date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	MachineID uint      `json:"machine_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MesinBor struct {
	gorm.Model
	Name              string  `gorm:"not null;unique" json:"name"`
	StockAvailability int     `gorm:"not null;default:0" json:"stock_availability"`
	RentalCosts       float64 `gorm:"not null;default:0" json:"rental_costs"`
	Category          string  `gorm:"default:'Uncategorized'" json:"category"`
	Description       string  `gorm:"size:255" json:"description"`
	Brand             string  `gorm:"size:100" json:"brand"`
	Condition         string  `gorm:"default:'Good';check:condition IN ('Good', 'Damaged', 'Needs Maintenance')" json:"condition"`
}

func GetMachines() ([]MesinBor, error) {
	var machines []MesinBor
	if err := DB.Find(&machines).Error; err != nil {
		return nil, err
	}
	return machines, nil
}

func GetMachineByID(id int) (*MesinBor, error) {
	var machine MesinBor
	if err := DB.First(&machine, id).Error; err != nil {
		return nil, err
	}
	return &machine, nil
}

func CreateMachine(machine *MesinBor) error {
	if err := DB.Create(machine).Error; err != nil {
		return err
	}
	return nil
}

func UpdateMachine(id int, updatedMachine *MesinBor) error {
	var machine MesinBor
	if err := DB.First(&machine, id).Error; err != nil {
		return err
	}

	if err := DB.Model(&machine).Updates(updatedMachine).Error; err != nil {
		return err
	}
	return nil
}

func DeleteMachine(id int) error {
	if err := DB.Delete(&MesinBor{}, id).Error; err != nil {
		return err
	}
	return nil
}

func CreateMaintenance(maintenance *Maintenance) error {
	if err := DB.Create(&maintenance).Error; err != nil {
		return err
	}
	return nil
}

func GetMaintenanceByID(id int) (*Maintenance, error) {
	var maintenance Maintenance
	if err := DB.First(&maintenance, id).Error; err != nil {
		return nil, err
	}
	return &maintenance, nil
}

func GetAllMaintenance() ([]Maintenance, error) {
	var records []Maintenance
	if err := DB.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func CreateRental(rental *RentalHistory) error {
	if err := DB.Create(&rental).Error; err != nil {
		return err
	}
	return nil
}

func GetRentalByID(id int) (*RentalHistory, error) {
	var rental RentalHistory
	if err := DB.First(&rental, id).Error; err != nil {
		return nil, err
	}
	return &rental, nil
}

func GetAllRentals() ([]RentalHistory, error) {
	var rentals []RentalHistory
	if err := DB.Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

func MarkAsReturned(id int, returnDate *time.Time) error {
	var rental RentalHistory
	if err := DB.First(&rental, id).Error; err != nil {
		return err
	}

	rental.ReturnDate = returnDate
	if err := DB.Save(&rental).Error; err != nil {
		return err
	}
	return nil
}

func CreateReview(review *Review) error {
	if err := DB.Create(&review).Error; err != nil {
		return err
	}
	return nil
}

func GetReviewByID(id int) (*Review, error) {
	var review Review
	if err := DB.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func GetAllReviews() ([]Review, error) {
	var reviews []Review
	if err := DB.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func DeleteReview(id int) error {
	if err := DB.Delete(&Review{}, id).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(user *User) error {
	if err := DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(id int) (*User, error) {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id int, user *User) error {
	var existingUser User
	if err := DB.First(&existingUser, id).Error; err != nil {
		return err
	}

	existingUser.Email = user.Email
	existingUser.Password = user.Password
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName

	if err := DB.Save(&existingUser).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(id int) error {
	if err := DB.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func CheckUserExists(email string) (bool, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func AuthenticateUser(email, password string) (*User, error) {
	var user User
	if err := DB.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	return &user, nil
}
