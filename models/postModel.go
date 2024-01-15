package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;"`
	CreatedAt time.Time
	Author    uuid.UUID
	Text      string
	Files     []File `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	post.ID = uuid.New()
	return
}
