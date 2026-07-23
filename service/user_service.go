package service

import (
    "spotsync/dto"
    "spotsync/repository"
)

type UserService interface {
    GetAllUsers() ([]dto.UserResponse, error)
    GetUserByID(id uint) (*dto.UserResponse, error)
    DeleteUser(id uint) error
}

type userServiceImpl struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) GetAllUsers() ([]dto.UserResponse, error) {
    users, err := s.userRepo.GetAllUsers()
    if err != nil {
        return nil, err
    }

    var resp []dto.UserResponse
    for _, u := range users {
        resp = append(resp, dto.UserResponse{
            ID:        u.ID,
            Name:      u.Name,
            Email:     u.Email,
            Role:      u.Role,
            CreatedAt: u.CreatedAt,
            UpdatedAt: u.UpdatedAt,
        })
    }

    return resp, nil
}

func (s *userServiceImpl) GetUserByID(id uint) (*dto.UserResponse, error) {
    u, err := s.userRepo.GetUserByID(id)
    if err != nil {
        return nil, err
    }

    return &dto.UserResponse{
        ID:        u.ID,
        Name:      u.Name,
        Email:     u.Email,
        Role:      u.Role,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }, nil
}

func (s *userServiceImpl) DeleteUser(id uint) error {
    return s.userRepo.DeleteUser(id)
}
