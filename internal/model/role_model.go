package model

import "time"

type RoleResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateRoleRequest struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required,max=100"`
}

type UpdateRoleRequest struct {
	ID   string `json:"-" validate:"required,max=100,uuid"`
	Name string `json:"name" validate:"required,max=100"`
}

type SearchRoleRequest struct {
	Name string `json:"name" validate:"max=100"`
	Page int    `json:"page" validate:"min=1"`
	Size int    `json:"size" validate:"min=1,max=100"`
}

type GetRoleRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}

type DeleteRoleRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}
