package models

import "github.com/google/uuid"

type File struct {
	ID       uint `gorm:"primaryKey;autoIncrement;"`
	FileName string
	PostID   *uuid.UUID
	UserID   *uuid.UUID
}
