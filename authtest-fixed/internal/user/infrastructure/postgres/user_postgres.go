package postgres

import (
	"database/sql"

	"github.com/test/authtest-fixed/internal/user/domain/entity"
	"github.com/test/authtest-fixed/internal/user/domain/repository"

)

// userRepository implements the UserRepository interface using PostgreSQL
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *entity.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, email, password_hash, is_active) VALUES ($1, $2, $3, $4)",
		user.Username, user.Email, user.PasswordHash, user.IsActive)
	return err
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users WHERE email = $1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *entity.User) error {
	_, err := r.db.Exec("UPDATE users SET username = $1, email = $2, password_hash = $3, is_active = $4, updated_at = NOW() WHERE id = $5",
		user.Username, user.Email, user.PasswordHash, user.IsActive, user.ID)
	return err
}

// Delete deletes a user
func (r *userRepository) Delete(id uint) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// FindAll finds all users with pagination
func (r *userRepository) FindAll(limit, offset int) ([]*entity.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// Count returns the total number of users
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// FindByIDWithRoles finds a user by ID with roles preloaded
func (r *userRepository) FindByIDWithRoles(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AssignRole assigns a role to a user
func (r *userRepository) AssignRole(userID, roleID uint) error {
	// Check if association already exists
	var count int64
	r.db.Table("user_roles").Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	if count > 0 {
		return nil // Already assigned
	}

	// Create association
	return r.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error
}

// RemoveRole removes a role from a user
func (r *userRepository) RemoveRole(userID, roleID uint) error {
	return r.db.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Error
}

// FindUsersByRole finds users by role name
func (r *userRepository) FindUsersByRole(roleName string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.name = ?", roleName).
		Find(&users).Error
	return users, err
}

// HasPermission checks if user has specific permission
func (r *userRepository) HasPermission(userID uint, permissionName string) (bool, error) {
	var count int64
	err := r.db.Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.id = ? AND permissions.name = ? AND users.is_active = ? AND roles.is_active = ? AND permissions.is_active = ?", 
			userID, permissionName, true, true, true).
		Count(&count).Error
	
	return count > 0, err
}

// GetUserPermissions gets all permissions for a user
func (r *userRepository) GetUserPermissions(userID uint) ([]string, error) {
	var permissions []string
	err := r.db.Table("permissions").
		Select("DISTINCT permissions.name").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.is_active = ? AND roles.is_active = ?", 
			userID, true, true).
		Pluck("permissions.name", &permissions).Error
	
	return permissions, err
}

// SearchByEmailOrUsername searches users by email or username
func (r *userRepository) SearchByEmailOrUsername(query string, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	searchPattern := "%" + query + "%"
	err := r.db.Where("email ILIKE ? OR username ILIKE ?", searchPattern, searchPattern).
		Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// FindUsersCreatedBetween finds users created between dates
func (r *userRepository) FindUsersCreatedBetween(startDate, endDate string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&users).Error
	return users, err
}
