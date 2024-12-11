package user

import (
	"finance-tracker/pkg/database"
)

type UserRepositoryInterface interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) (*User, error)
}

type UserRepository struct {
	Database *database.Database
}

func NewUserRepository(database *database.Database) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
