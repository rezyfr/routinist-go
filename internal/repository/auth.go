package repository

import (
	"math/rand"
	"os"
	"routinist/internal/domain/errors"
	"routinist/internal/domain/model"
	"routinist/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewAuthRepo(db *gorm.DB, logger *logger.Logger) *authRepo {
	return &authRepo{
		db, logger,
	}
}

func (rp *authRepo) Register(e *model.RegisterRequestDTO) (*model.AuthResponseDTO, error) {
	// Check if email already exists
	var user model.User

	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)
	exists := result.RowsAffected > 0
	if exists {
		return &model.AuthResponseDTO{}, errors.ErrEmailAlreadyExists
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return &model.AuthResponseDTO{}, errors.ErrFailedToHashPassword
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
	}
	result = rp.db.Create(&user)
	if result.Error != nil {
		return &model.AuthResponseDTO{}, result.Error
	}

	token, err := generateJWT(user.Email)
	if err != nil {
		return &model.AuthResponseDTO{}, errors.ErrFailedToGenerateJWT
	}

	return &model.AuthResponseDTO{Token: token}, nil
}

func (rp *authRepo) Login(e *model.LoginRequestDTO) (*model.AuthResponseDTO, error) {
	var user model.User
	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)

	if result.Error != nil {
		return &model.AuthResponseDTO{}, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(e.Password)); err != nil {
		return &model.AuthResponseDTO{}, errors.ErrInvalidCredentials
	}

	token, err := generateJWT(user.Email)
	if err != nil {
		return &model.AuthResponseDTO{}, errors.ErrFailedToGenerateJWT
	}

	return &model.AuthResponseDTO{Token: token}, nil
}

func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
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
