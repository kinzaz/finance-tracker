package auth

import (
	"errors"
	"finance-tracker/internal/errs"
	"finance-tracker/internal/models"
	"finance-tracker/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(dto *RegisterRequestDto) (*RegisterResponseDto, error)
	Login(dto *LoginRequestDto) (*models.User, error)
}

type AuthService struct {
	UserRepository user.UserRepositoryInterface
}

func NewAuthService(userRepository user.UserRepositoryInterface) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Login(dto *LoginRequestDto) (*models.User, error) {
	existedUser, _ := service.UserRepository.FindByEmail(dto.Email)

	if existedUser == nil {
		return existedUser, errors.New("неверный логин или пароль")
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(dto.Password))
	if err != nil {
		return existedUser, errors.New("неверный логин или пароль")
	}

	return existedUser, nil
}

func (service *AuthService) Register(dto *RegisterRequestDto) (*RegisterResponseDto, error) {
	existedUser, _ := service.UserRepository.FindByEmail(dto.Email)
	if existedUser != nil {
		return nil, errs.ErrRegisterUser
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userEntity := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: string(hashedPassword),
	}

	_, err = service.UserRepository.Create(userEntity)
	if err != nil {
		return nil, err
	}

	response := &RegisterResponseDto{
		ID: userEntity.ID,
	}

	return response, nil
}
