package repository

import (
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

func (rp *AuthRepo) Register(db *gorm.DB, e *request.RegisterRequestDTO) (*request.AuthResponseDTO, uint, error) {
	// Check if email already exists
	var user model.User

	result := db.Where("email = ?", e.Email).Limit(1).Find(&user)
	exists := result.RowsAffected > 0
	if exists {
		return &request.AuthResponseDTO{}, 0, errors.ErrEmailAlreadyExists
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return &request.AuthResponseDTO{}, 0, errors.ErrFailedToHashPassword
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

	result = db.Create(&user)
	if result.Error != nil {
		return &request.AuthResponseDTO{}, 0, result.Error
	}

	rp.logger.Info("User created: ", user)

	token, err := auth.GenerateJWT(user.Email, user.ID)
	if err != nil {
		return &request.AuthResponseDTO{}, 0, errors.ErrFailedToGenerateJWT
	}

	return &request.AuthResponseDTO{Token: token}, user.ID, nil
}

func (rp *AuthRepo) Login(e *request.LoginRequestDTO) (*request.AuthResponseDTO, uint, error) {
	var user model.User
	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)

	if result.Error != nil {
		return &request.AuthResponseDTO{}, 0, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(e.Password)); err != nil {
		return &request.AuthResponseDTO{}, 0, errors.ErrInvalidCredentials
	}

	token, err := auth.GenerateJWT(user.Email, user.ID)
	if err != nil {
		return &request.AuthResponseDTO{}, 0, errors.ErrFailedToGenerateJWT
	}

	return &request.AuthResponseDTO{Token: token}, user.ID, nil
}

func (rp *AuthRepo) GetDB() *gorm.DB {
	return rp.db
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
