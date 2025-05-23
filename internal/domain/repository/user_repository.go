package repository

type UserRepository interface {
	UpdateMilestone(userId uint, milestone uint) (uint, error)
}
