package converter

import (
	"todo-app/internal/entity"
	"todo-app/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {

	return &model.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		RoleId:        user.RoleId,
		DepartementId: user.DepartementId,
		IsActive:      user.IsActive,
	}
}

func UserToTokenResponse(user *entity.User, expiresIn int64) *model.UserResponse {
	return &model.UserResponse{
		Token:     user.Token,
		ExpiresIn: expiresIn,
	}
}
