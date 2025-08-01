package entity

import "time"

// Permission represents a permission in the RBAC system
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	Resource    string    `json:"resource" gorm:"not null"` // e.g., "users", "products", "orders"
	Action      string    `json:"action" gorm:"not null"`   // e.g., "create", "read", "update", "delete"
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Note: RBAC relationships are managed through repository layer
	// to avoid circular imports between bounded contexts
}

// PermissionConstants defines common permission constants
var PermissionConstants = struct {
	// User permissions
	UserCreate string
	UserRead   string
	UserUpdate string
	UserDelete string
	
	// Role permissions
	RoleCreate string
	RoleRead   string
	RoleUpdate string
	RoleDelete string
	
	// Permission permissions
	PermissionCreate string
	PermissionRead   string
	PermissionUpdate string
	PermissionDelete string
	
	// Admin permissions
	AdminAll string
}{
	// User permissions
	UserCreate: "user:create",
	UserRead:   "user:read",
	UserUpdate: "user:update",
	UserDelete: "user:delete",
	
	// Role permissions
	RoleCreate: "role:create",
	RoleRead:   "role:read",
	RoleUpdate: "role:update",
	RoleDelete: "role:delete",
	
	// Permission permissions
	PermissionCreate: "permission:create",
	PermissionRead:   "permission:read",
	PermissionUpdate: "permission:update",
	PermissionDelete: "permission:delete",
	
	// Admin permissions
	AdminAll: "admin:all",
}

// GetFullName returns the full permission name in format "resource:action"
func (p *Permission) GetFullName() string {
	return p.Resource + ":" + p.Action
}

// IsResourceAction checks if permission matches specific resource and action
func (p *Permission) IsResourceAction(resource, action string) bool {
	return p.Resource == resource && p.Action == action
}
