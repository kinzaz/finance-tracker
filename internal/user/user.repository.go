package user

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/internal/models"
	"finance-tracker/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindByEmail(email string) (*models.User, error)
	FindById(id uint) (*models.User, error)
	Create(user *models.User) (*models.User, error)
}

type UserRepository struct {
	Database *database.Database
}

func NewUserRepository(database *database.Database) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindById(id uint) (*models.User, error) {
	var user models.User
	result := repo.Database.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: user_id = %d", errs.ErrUserNotFound, id)
		}
		return nil, result.Error
	}
	return &user, nil
}

func (repo UserRepository) Create(user *models.User) (*models.User, error) {
	result := repo.Database.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
