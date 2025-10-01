package model

import "time"

type ContactResponse struct {
	ID        string            `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt time.Time         `json:"deleted_at"`
	Addresses []AddressResponse `json:"addresses,omitempty"`
}

type CreateContactRequest struct {
	ID        string `json:"id"`
	UserId    string `json:"-" validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
	Email     string `json:"email" validate:"max=200,email"`
	Phone     string `json:"phone" validate:"max=20"`
}

type UpdateContactRequest struct {
	UserId    string `json:"-" validate:"required"`
	ID        string `json:"-" validate:"required,max=100,uuid"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
	Email     string `json:"email" validate:"max=200,email"`
	Phone     string `json:"phone" validate:"max=20"`
}

type SearchContactRequest struct {
	UserId string `json:"-" validate:"required"`
	Name   string `json:"name" validate:"max=100"`
	Email  string `json:"email" validate:"max=200"`
	Phone  string `json:"phone" validate:"max=20"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}

type GetContactRequest struct {
	UserId string `json:"-" validate:"required"`
	ID     string `json:"-" validate:"required,max=100,uuid"`
}

type DeleteContactRequest struct {
	UserId string `json:"-" validate:"required"`
	ID     string `json:"-" validate:"required,max=100,uuid"`
}
