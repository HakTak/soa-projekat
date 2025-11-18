package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleGuide   Role = "GUIDE"
	RoleTourist Role = "TOURIST"
	RoleAdmin   Role = "ADMIN"
)

type User struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Username string    `json:"username" gorm:"not null;unique"`
	Password string    `json:"password"`
	Email    string    `json:"email" gorm:"unique"`
	Role     Role      `json:"role"`
	Blocked  bool      `json:"blocked" gorm:"default:false"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	return
}
