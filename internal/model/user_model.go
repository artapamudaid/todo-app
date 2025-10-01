package model

type UserResponse struct {
	ID            string `json:"id,omitempty"`
	Email         string `json:"email,omitempty"`
	Name          string `json:"name,omitempty"`
	RoleId        string `json:"role_id,omitempty"`
	DepartementId string `json:"department_id,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
	Token         string `json:"token,omitempty"`
	ExpiresIn     int64  `json:"expires_in,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=500"`
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	ID       string `json:"-" validate:"required,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Password string `json:"password,omitempty" validate:"max=100"`
	Name     string `json:"name,omitempty" validate:"max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
