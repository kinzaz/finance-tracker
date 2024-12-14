package user

import "finance-tracker/internal/models"

type UserServiceInterface interface {
	GetUserProfile(id uint) (*models.User, error)
}

type UserService struct {
	UserRepository UserRepositoryInterface
}

func NewUserService(userRepository UserRepositoryInterface) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (service *UserService) GetUserProfile(id uint) (*models.User, error) {
	userProfile, err := service.UserRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}
