package entity

import (
	"time"

	"gorm.io/gorm"
)

// User is a struct that represents a user entity
type User struct {
	ID            string         `gorm:"column:id;primaryKey"`
	Email         string         `gorm:"column:email;uniqueIndex"`
	Password      string         `gorm:"column:password"`
	Name          string         `gorm:"column:name"`
	RoleId        string         `gorm:"column:role_id"`
	DepartementId string         `gorm:"column:department_id"`
	IsActive      bool           `gorm:"column:is_active"`
	Token         string         `gorm:"column:token"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (u *User) TableName() string {
	return "users"
}
