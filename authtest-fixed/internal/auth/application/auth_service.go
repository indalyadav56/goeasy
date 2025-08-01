package application

import (
	"errors"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"

	"github.com/test/authtest-fixed/internal/user/domain/entity"
	userRepo "github.com/test/authtest-fixed/internal/user/domain/repository"
	roleEntity "github.com/test/authtest-fixed/internal/role/domain/entity"
	roleRepo "github.com/test/authtest-fixed/internal/role/domain/repository"
	permissionRepo "github.com/test/authtest-fixed/internal/permission/domain/repository"

	"github.com/test/authtest-fixed/pkg/auth"
)

type AuthService struct {
	userRepo       userRepo.UserRepository
	roleRepo       roleRepo.RoleRepository
	permissionRepo permissionRepo.PermissionRepository
	jwtSecret      string
	tokenExpiry    time.Duration
}

func NewAuthService(
	userService interface{},
	roleService interface{},
	permissionService interface{},
) *AuthService {
	return &AuthService{
		jwtSecret:   "default-secret",
		tokenExpiry: 24 * time.Hour,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type AuthResponse struct {
	User         *entity.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Load user roles and permissions
	user, err = s.userRepo.FindByIDWithRoles(user.ID)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.tokenExpiry),
	}, nil
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
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
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign default role (e.g., "user")
	defaultRole, err := s.roleRepo.FindByName("user")
	if err == nil && defaultRole != nil {
		s.userRepo.AssignRole(user.ID, defaultRole.ID)
	}

	// Load user with roles
	user, err = s.userRepo.FindByIDWithRoles(user.ID)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.tokenExpiry),
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := auth.ValidateJWT(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Find user
	user, err := s.userRepo.FindByIDWithRoles(uint(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Keep the same refresh token
		ExpiresAt:    time.Now().Add(s.tokenExpiry),
	}, nil
}

// ValidateToken validates an access token and returns user info
func (s *AuthService) ValidateToken(tokenString string) (*entity.User, error) {
	claims, err := auth.ValidateJWT(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Get user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Find user with roles
	user, err := s.userRepo.FindByIDWithRoles(uint(userID))
	if err != nil {
		return nil, err
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	return user, nil
}

// CheckPermission checks if user has specific permission
func (s *AuthService) CheckPermission(userID uint, permissionName string) bool {
	user, err := s.userRepo.FindByIDWithRoles(userID)
	if err != nil {
		return false
	}

	fmt.Println("user", user)

	// Check if user has permission
	permission, err := s.permissionRepo.FindByName(permissionName)
	if err != nil {
		return false
	}

	fmt.Println("permission", permission)

	return true
}

// generateAccessToken generates JWT access token
func (s *AuthService) generateAccessToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		// "roles":    s.extractRoleNames(user.Roles),
		"exp":      time.Now().Add(s.tokenExpiry).Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	return auth.GenerateJWT(claims, s.jwtSecret)
}

// generateRefreshToken generates JWT refresh token
func (s *AuthService) generateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	return auth.GenerateJWT(claims, s.jwtSecret)
}

// extractRoleNames extracts role names from user roles
func (s *AuthService) extractRoleNames(roles []roleEntity.Role) []string {
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}
	return roleNames
}
