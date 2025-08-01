package handlers

import (
	"net/http"
	"strconv"
	
	"github.com/test/authtest-fixed/internal/permission/application"
	"github.com/test/authtest-fixed/internal/permission/interface/http/v1/dto"
	"github.com/go-chi/chi/v5"
	"encoding/json"
)

type PermissionHandler struct {
	permissionService *application.PermissionService
}

func NewPermissionHandler(permissionService *application.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}
// GetPermissions returns a list of permissions
func (h *PermissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 10
	}
	
	permissions, total, err := h.permissionService.GetPermissions(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"permissions":  permissions,
		"total_count":  total,
		"page":         page,
		"page_size":    pageSize,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPermissionByID returns a permission by ID
func (h *PermissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	
	permission, err := h.permissionService.GetPermissionByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, "Permission not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permission)
}

// CreatePermission creates a new permission
func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	permission, err := h.permissionService.CreatePermission(r.Context(), req.Name, req.Resource, req.Action, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(permission)
}

// UpdatePermission updates a permission
func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	
	var req dto.UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	permission, err := h.permissionService.UpdatePermission(r.Context(), uint(id), req.Name, req.Resource, req.Action, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permission)
}

// DeletePermission deletes a permission
func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	
	err = h.permissionService.DeletePermission(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission deleted successfully"})
}

// SearchPermissions searches permissions by name or resource
func (h *PermissionHandler) SearchPermissions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}
	
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 10
	}
	
	permissions, total, err := h.permissionService.SearchPermissions(r.Context(), query, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"permissions":  permissions,
		"total_count":  total,
		"page":         page,
		"page_size":    pageSize,
		"query":        query,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
