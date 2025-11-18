package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role is a string alias for user roles.
type Role string

const (
	RoleGuide   Role = "guide"
	RoleTourist Role = "tourist"
	RoleAdmin   Role = "admin"
)

type Profile struct {
	ID             string `gorm:"primaryKey;type:uuid" json:"id"`
	UserID         string `gorm:"uniqueIndex;type:uuid;not null" json:"user_id"`
	FirstName      string `gorm:"type:text;not null" json:"first_name"`
	LastName       string `gorm:"type:text;not null" json:"last_name"`
	ProfilePicture string `gorm:"type:text" json:"profile_picture"`
	Biography      string `gorm:"type:text" json:"biography"`
	Motto          string `gorm:"type:text" json:"motto"`

	Role      Role `gorm:"type:varchar(20);not null" json:"role"`
	IsBlocked bool `gorm:"type:boolean;default:false;not null" json:"is_blocked"`
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}
