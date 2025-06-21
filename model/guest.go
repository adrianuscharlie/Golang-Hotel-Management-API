package model

import "time"

type Guest struct {
	ID             uint `gorm:"primaryKey"`
	CredentialType string
	FullName       string
	Phone          string
	Email          string
	IDNumber       string `gorm:"unique"`
	CreatedAt      time.Time
}
