package entity

import "time"

// Role represents a role in the RBAC system
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Description string       `json:"description"`
	IsActive    bool         `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	
	// Note: RBAC relationships are managed through repository layer
	// to avoid circular imports between bounded contexts
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uint      `json:"role_id" gorm:"primaryKey"`
	PermissionID uint      `json:"permission_id" gorm:"primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
}

// Note: Role-Permission relationship methods are handled through repository layer
// to avoid circular imports between bounded contexts
