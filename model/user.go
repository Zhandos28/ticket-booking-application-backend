package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleManager  UserRole = "manager"
	UserRoleAttendee UserRole = "attendee"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;type:integer;autoIncrement"`
	Email     string    `json:"email" gorm:"type:text;unique;not null"`
	Role      UserRole  `json:"role" gorm:"type:text;default:attendee"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", UserRoleManager)
	}
	return
}
