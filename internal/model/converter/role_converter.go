package converter

import (
	"todo-app/internal/entity"
	"todo-app/internal/model"
)

func RoleToResponse(role *entity.Role) *model.RoleResponse {
	return &model.RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
