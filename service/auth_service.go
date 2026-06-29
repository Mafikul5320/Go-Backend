package service

import (
	"errors"

	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
	"spotsync/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (string, *dto.UserResponse, error)
}

type authServiceImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authServiceImpl{userRepo: userRepo}
}

func (s *authServiceImpl) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email exists
	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	role := "driver"
	if req.Role == "admin" {
		role = "admin"
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *authServiceImpl) Login(req dto.LoginRequest) (string, *dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
