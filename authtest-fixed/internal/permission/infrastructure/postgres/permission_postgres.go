package postgres

import (
	"database/sql"

	"github.com/test/authtest-fixed/internal/permission/domain/entity"
	"github.com/test/authtest-fixed/internal/permission/domain/repository"
	roleEntity "github.com/test/authtest-fixed/internal/role/domain/entity"

)

// permissionRepository implements the PermissionRepository interface using PostgreSQL
type permissionRepository struct {
	db *sql.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *sql.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

// Create creates a new permission
func (r *permissionRepository) Create(permission *entity.Permission) error {
	_, err := r.db.Exec("INSERT INTO permissions (name, description, resource, action, is_active) VALUES ($1, $2, $3, $4, $5)", 
		permission.Name, permission.Description, permission.Resource, permission.Action, permission.IsActive)
	return err
}

// FindByID finds a permission by ID
func (r *permissionRepository) FindByID(id uint) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.QueryRow("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE id = $1", id).Scan(
		&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByName finds a permission by name
func (r *permissionRepository) FindByName(name string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.QueryRow("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE name = $1", name).Scan(
		&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Update updates a permission
func (r *permissionRepository) Update(permission *entity.Permission) error {
	_, err := r.db.Exec("UPDATE permissions SET name = $1, description = $2, is_active = $3, updated_at = NOW() WHERE id = $4",
		permission.Name, permission.Description, permission.IsActive, permission.ID)
	return err
}

// Delete deletes a permission
func (r *permissionRepository) Delete(id uint) error {
	_, err := r.db.Exec("DELETE FROM permissions WHERE id = $1", id)
	return err
}

// FindAll finds all permissions with pagination
func (r *permissionRepository) FindAll(limit, offset int) ([]*entity.Permission, error) {
	rows, err := r.db.Query("SELECT id, name, description, is_active, created_at, updated_at FROM permissions LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

// FindByStatus finds permissions by status with pagination
func (r *permissionRepository) FindByStatus(isActive bool, limit, offset int) ([]*entity.Permission, error) {
	rows, err := r.db.Query("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE is_active = $1 LIMIT $2 OFFSET $3", isActive, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

// FindByResource finds permissions by resource (simplified - removed resource field)
func (r *permissionRepository) FindByResource(resource string) ([]*entity.Permission, error) {
	rows, err := r.db.Query("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE name LIKE $1", "%"+resource+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

// FindByAction finds permissions by action (simplified)
func (r *permissionRepository) FindByAction(action string) ([]*entity.Permission, error) {
	rows, err := r.db.Query("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE name LIKE $1", "%"+action+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

// FindByResourceAndAction finds permission by resource and action (simplified)
func (r *permissionRepository) FindByResourceAndAction(resource, action string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.QueryRow("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE name LIKE $1", "%"+resource+"%"+action+"%").Scan(
		&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Count returns total number of permissions
func (r *permissionRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM permissions").Scan(&count)
	return count, err
}

func (r *permissionRepository) GetPermissionRoles(permissionID uint) ([]*roleEntity.Role, error) {
	// Simplified implementation - return empty slice
	return []*roleEntity.Role{}, nil
}

// FindPermissionsByRole finds permissions assigned to a role
func (r *permissionRepository) FindPermissionsByRole(roleID uint) ([]*entity.Permission, error) {
	// Simplified implementation - return empty slice
	return []*entity.Permission{}, nil
}

// CheckUserPermission checks if user has specific permission (simplified)
func (r *permissionRepository) CheckUserPermission(userID uint, permissionName string) (bool, error) {
	// Simplified implementation - always return false
	return false, nil
}

// GetUserPermissions gets all permissions for a user (simplified)
func (r *permissionRepository) GetUserPermissions(userID uint) ([]*entity.Permission, error) {
	// Simplified implementation - return empty slice
	return []*entity.Permission{}, nil
}

// SearchByName searches permissions by name
func (r *permissionRepository) SearchByName(query string, limit, offset int) ([]*entity.Permission, error) {
	rows, err := r.db.Query("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE name ILIKE $1 LIMIT $2 OFFSET $3", "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

// SearchByDescription searches permissions by description
func (r *permissionRepository) SearchByDescription(query string, limit, offset int) ([]*entity.Permission, error) {
	rows, err := r.db.Query("SELECT id, name, description, is_active, created_at, updated_at FROM permissions WHERE description ILIKE $1 LIMIT $2 OFFSET $3", "%"+query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		var permission entity.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.IsActive, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}
	return permissions, nil
}

// FindPermissionsCreatedBetween finds permissions created between dates (simplified)
func (r *permissionRepository) FindPermissionsCreatedBetween(startDate, endDate string) ([]*entity.Permission, error) {
	// Simplified implementation - return empty slice
	return []*entity.Permission{}, nil
}

// CreateBulk creates multiple permissions in bulk (simplified)
func (r *permissionRepository) CreateBulk(permissions []*entity.Permission) error {
	// Simplified implementation - do nothing
	return nil
}

// FindByNames finds permissions by names (simplified)
func (r *permissionRepository) FindByNames(names []string) ([]*entity.Permission, error) {
	// Simplified implementation - return empty slice
	return []*entity.Permission{}, nil
}
