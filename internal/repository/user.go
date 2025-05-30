package repository

import (
	"gorm.io/gorm"
	"routinist/internal/domain/model"
	"routinist/pkg/logger"
)

type UserRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewUserRepo(db *gorm.DB, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		db, logger,
	}
}

func (rp *UserRepo) UpdateMilestone(userId uint, milestone uint) (uint, error) {
	var user model.User
	err := rp.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return 0, err
	}

	rp.logger.Info("User milestone before update: ", user)
	user.Milestone = user.Milestone + milestone
	rp.logger.Info("User milestone updated: ", user)
	err = rp.db.Save(&user).Error
	if err != nil {
		return 0, err
	}

	return user.Milestone, nil
}
