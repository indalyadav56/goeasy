package repository

import (

	"github.com/test/authtest-fixed/internal/role/domain/entity"
	userEntity "github.com/test/authtest-fixed/internal/user/domain/entity"

)

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	// Basic CRUD operations
	Create(role *entity.Role) error
	FindByID(id uint) (*entity.Role, error)
	FindByName(name string) (*entity.Role, error)
	Update(role *entity.Role) error
	Delete(id uint) error
	
	// Role listing and filtering
	FindAll(limit, offset int) ([]*entity.Role, error)
	FindByStatus(isActive bool, limit, offset int) ([]*entity.Role, error)
	Count() (int64, error)
	
	// Permission management
	FindByIDWithPermissions(id uint) (*entity.Role, error)
	AssignPermission(roleID, permissionID uint) error
	RemovePermission(roleID, permissionID uint) error
	FindRolesByPermission(permissionName string) ([]*entity.Role, error)
	
	// User management
	GetRoleUsers(roleID uint) ([]*userEntity.User, error)
	GetUserRoles(userID uint) ([]*entity.Role, error)
	
	// Search and filtering
	SearchByName(query string, limit, offset int) ([]*entity.Role, error)
	FindRolesCreatedBetween(startDate, endDate string) ([]*entity.Role, error)
}
