package application

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"github.com/test/authtest-fixed/internal/user/domain/entity"
	"github.com/test/authtest-fixed/internal/user/domain/repository"
	roleEntity "github.com/test/authtest-fixed/internal/role/domain/entity"
	roleRepo "github.com/test/authtest-fixed/internal/role/domain/repository"

)

type UserService struct {
	userRepo repository.UserRepository
	roleRepo roleRepo.RoleRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		// roleRepo: nil, // Simplified - not injected
	}
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"email"`
	Username  string `json:"username" validate:"min=3"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  *bool  `json:"is_active"`
}

func (s *UserService) GetUsers(limit, offset int) ([]*entity.User, error) {
	return s.userRepo.FindAll(limit, offset)
}

func (s *UserService) GetUser(id uint) (*entity.User, error) {
	return s.userRepo.FindByIDWithRoles(id)
}

func (s *UserService) CreateUser(req CreateUserRequest) (*entity.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUser, _ = s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Create new user
	user := &entity.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// Hash password
	err := user.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(id uint, req UpdateUserRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Email != "" && req.Email != user.Email {
		// Check if email is already taken
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already taken")
		}
		user.Email = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		// Check if username is already taken
		existingUser, _ := s.userRepo.FindByUsername(req.Username)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("username already taken")
		}
		user.Username = req.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.userRepo.FindByIDWithRoles(id)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(user.ID)
}

// AssignRole assigns a role to a user
func (s *UserService) AssignRole(userID, roleID uint) error {
	// Verify user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify role exists
	_, err = s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	return s.userRepo.AssignRole(userID, roleID)
}

// RemoveRole removes a role from a user
func (s *UserService) RemoveRole(userID, roleID uint) error {
	// Verify user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.RemoveRole(userID, roleID)
}

func (s *UserService) GetUserRoles(userID uint) ([]*roleEntity.Role, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.roleRepo.GetUserRoles(userID)
}

func (s *UserService) SearchUsers(query string, limit, offset int) ([]*entity.User, error) {
	return s.userRepo.SearchByEmailOrUsername(query, limit, offset)
}

func (s *UserService) GetUserCount() (int64, error) {
	return s.userRepo.Count()
}

func (s *UserService) CheckUserPermission(userID uint, permissionName string) (bool, error) {
	return s.userRepo.HasPermission(userID, permissionName)
}
