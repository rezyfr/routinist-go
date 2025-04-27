package repo

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"routinist/internal/entity"
	"routinist/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
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

func (rp *AuthRepo) Register(ctx context.Context, e entity.RegisterRequestDTO) (entity.AuthResponseDTO, error) {
	// Check if email already exists
	var user entity.User

	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)
	exists := result.RowsAffected > 0
	if exists {
		return entity.AuthResponseDTO{}, fmt.Errorf("email already exists")
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.AuthResponseDTO{}, err
	}

	// Set default name if not provided
	name := e.Name
	if name == "" {
		name = generateRandomName()
	}

	// Create user in database
	user = entity.User{
		Email:    e.Email,
		Password: string(hash),
		Name:     name,
	}
	result = rp.db.Create(&user)
	if result.Error != nil {
		return entity.AuthResponseDTO{}, result.Error
	}

	token, err := generateJWT(user.Email)
	if err != nil {
		return entity.AuthResponseDTO{}, err
	}

	return entity.AuthResponseDTO{Token: token}, nil
}

func (rp *AuthRepo) Login(ctx context.Context, e entity.LoginRequestDTO) (entity.AuthResponseDTO, error) {
	var user entity.User
	result := rp.db.Where("email = ?", e.Email).Limit(1).Find(&user)

	if result.Error != nil {
		return entity.AuthResponseDTO{}, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(e.Password)); err != nil {
		return entity.AuthResponseDTO{}, err
	}

	token, err := generateJWT(user.Email)
	if err != nil {
		return entity.AuthResponseDTO{}, err
	}

	return entity.AuthResponseDTO{Token: token}, nil
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
