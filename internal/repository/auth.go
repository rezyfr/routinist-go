package repository

import (
	"fmt"
	"math/rand"
	"routinist/internal/auth"
	"routinist/internal/domain/errors"
	"routinist/internal/domain/model"
	"routinist/internal/dto/request"
	"routinist/pkg/logger"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewAuthRepo(db *gorm.DB, logger *logger.Logger) *AuthRepo {
	return &AuthRepo{
		db, logger,
	}
}

func (rp *AuthRepo) Register(e *request.RegisterRequestDTO) (*request.AuthResponseDTO, error) {
	// Check if email already exists
	var user model.User

	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)
	exists := result.RowsAffected > 0
	if exists {
		return &request.AuthResponseDTO{}, errors.ErrEmailAlreadyExists
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return &request.AuthResponseDTO{}, errors.ErrFailedToHashPassword
	}

	// Set default name if not provided
	name := e.Name
	if name == "" {
		name = generateRandomName()
	}

	// Create user in database
	user = model.User{
		Email:    e.Email,
		Password: string(hash),
		Name:     name,
		Gender:   e.Gender,
	}

	result = rp.db.Create(&user)
	if result.Error != nil {
		return &request.AuthResponseDTO{}, result.Error
	}

	rp.logger.Info("User created: ", user)

	// Check if habit valid
	var habit model.Habit
	habitResult := rp.db.Preload("Units").Where("id = ?", e.HabitID).First(&habit)
	habitExists := habitResult.RowsAffected > 0
	if !habitExists {
		return &request.AuthResponseDTO{}, fmt.Errorf("habit %d does not exist", e.HabitID)
	}

	rp.logger.Info("Habit found: ", habit.Name)

	unit := habit.Units[0]
	// Link habit to user
	userHabit := model.UserHabit{
		UserID:  user.ID,
		HabitID: habit.ID,
		UnitID:  unit.ID,
		Goal:    habit.DefaultGoal,
	}

	if err := rp.db.Create(&userHabit).Error; err != nil {
		rp.logger.Error("Failed to create user habit: ", err)
		return &request.AuthResponseDTO{}, errors.ErrFailedToAddHabit
	}

	rp.logger.Info("User habit created: ", userHabit)
	token, err := auth.GenerateJWT(user.Email, user.ID)
	if err != nil {
		return &request.AuthResponseDTO{}, errors.ErrFailedToGenerateJWT
	}

	return &request.AuthResponseDTO{Token: token}, nil
}

func (rp *AuthRepo) Login(e *request.LoginRequestDTO) (*request.AuthResponseDTO, error) {
	var user model.User
	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)

	if result.Error != nil {
		return &request.AuthResponseDTO{}, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(e.Password)); err != nil {
		return &request.AuthResponseDTO{}, errors.ErrInvalidCredentials
	}

	token, err := auth.GenerateJWT(user.Email, user.ID)
	if err != nil {
		return &request.AuthResponseDTO{}, errors.ErrFailedToGenerateJWT
	}

	return &request.AuthResponseDTO{Token: token}, nil
}

// Randomize name consisted of 2 words, 1. Color 2. Animal. Each 20
func generateRandomName() string {
	colors := []string{
		"Red", "Blue", "Green", "Yellow", "Purple",
		"Orange", "Pink", "Brown", "Gray", "White",
		"Black", "Cyan", "Magenta", "Teal", "Indigo",
		"Violet", "Crimson", "Azure", "Coral", "Gold",
	}

	animals := []string{
		"Lion", "Tiger", "Bear", "Wolf", "Fox",
		"Eagle", "Hawk", "Owl", "Deer", "Rabbit",
		"Dragon", "Phoenix", "Unicorn", "Griffin", "Panther",
		"Falcon", "Raven", "Snake", "Leopard", "Dolphin",
	}

	rand.Seed(time.Now().UnixNano())
	color := colors[rand.Intn(len(colors))]
	animal := animals[rand.Intn(len(animals))]

	return color + " " + animal
}
