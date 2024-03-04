package models

import "github.com/google/uuid"

type File struct {
	ID       uint `gorm:"primaryKey;autoIncrement;"`
	FileName string
	PostID   *uuid.UUID `json:"-"`
	UserID   *uuid.UUID `json:"-"`
}
