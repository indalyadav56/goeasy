package handlers

import (
	"net/http"
	
	"github.com/test/authtest-fixed/internal/user/application"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	// if page == 0 {
	// 	page = 1
	// }
	// pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	// if pageSize == 0 {
	// 	pageSize = 10
	// }
	
	// users, total, err := h.userService.GetUsers(r.Context(), page, pageSize)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// response := dto.PaginatedUsersResponse{
	// 	Users:      users,
	// 	TotalCount: total,
	// 	Page:       page,
	// 	PageSize:   pageSize,
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	
	// user, err := h.userService.GetUserByID(r.Context(), uint(id))
	// if err != nil {
	// 	http.Error(w, "User not found", http.StatusNotFound)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	
	// var req dto.UpdateUserRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// user, err := h.userService.UpdateUser(r.Context(), uint(id), req)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	
	// err = h.userService.DeleteUser(r.Context(), uint(id))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// idStr := chi.URLParam(r, "id")
	// id, err := strconv.ParseUint(idStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	
	// var req dto.ChangePasswordRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// err = h.userService.ChangePassword(r.Context(), uint(id), req.OldPassword, req.NewPassword)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	// var req dto.AssignRoleRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	
	// err := h.userService.AssignRole(r.Context(), req.UserID, req.RoleID)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Role assigned successfully"})
}

func (h *UserHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	// userIDStr := chi.URLParam(r, "user_id")
	// userID, err := strconv.ParseUint(userIDStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }
	
	// roleIDStr := chi.URLParam(r, "role_id")
	// roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid role ID", http.StatusBadRequest)
	// 	return
	// }
	
	// err = h.userService.RemoveRole(r.Context(), uint(userID), uint(roleID))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Role removed successfully"})
}
