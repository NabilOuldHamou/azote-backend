package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Avatar    File   `gorm:"foreignKey:UserID;"`
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Password  string `json:"-"`
	Posts     []Post `gorm:"foreignKey:Author;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
