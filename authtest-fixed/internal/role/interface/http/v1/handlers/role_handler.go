package handlers

import (
	"net/http"
	"strconv"
	
	"github.com/test/authtest-fixed/internal/role/application"
)

type RoleHandler struct {
	roleService *application.RoleService
}

func NewRoleHandler(roleService *application.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}
func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 10
	}
	
	// roles, err := h.roleService.GetRoles(r.Context(), page, pageSize)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// response := dto.PaginatedRolesResponse{
	// 	Roles:      roles,
	// 	TotalCount: total,
	// 	Page:       page,
	// 	PageSize:   pageSize,
	// }
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid role ID", http.StatusBadRequest)
	// 	return
	// }
	
	// role, err := h.roleService.GetRoleByID(r.Context(), uint(id))
	// if err != nil {
	// 	http.Error(w, "Role not found", http.StatusNotFound)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// var req dto.CreateRoleRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// role, err := h.roleService.CreateRole(r.Context(), req.Name, req.Description)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid role ID", http.StatusBadRequest)
	// 	return
	// }
	
	// var req dto.UpdateRoleRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// role, err := h.roleService.UpdateRole(r.Context(), uint(id), req.Name, req.Description)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid role ID", http.StatusBadRequest)
	// 	return
	// }
	
	// err = h.roleService.DeleteRole(r.Context(), uint(id))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Role deleted successfully"})
}

func (h *RoleHandler) AssignPermission(w http.ResponseWriter, r *http.Request) {
	// var req dto.AssignPermissionRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// err := h.roleService.AssignPermission(r.Context(), req.RoleID, uint(1))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Permission assigned successfully"})
}

func (h *RoleHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	// roleIDStr := chi.URLParam(r, "role_id")
	// roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid role ID", http.StatusBadRequest)
	// 	return
	// }
	
	// permissionIDStr := chi.URLParam(r, "permission_id")
	// permissionID, err := strconv.ParseUint(permissionIDStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid permission ID", http.StatusBadRequest)
	// 	return
	// }
	
	// err = h.roleService.RemovePermission(r.Context(), uint(roleID), uint(permissionID))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Permission removed successfully"})e(map[string]string{"message": "Permission removed successfully"})
}
