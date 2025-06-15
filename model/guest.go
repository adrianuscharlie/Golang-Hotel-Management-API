package model

type Guest struct {
	ID             string `gorm:"primaryKey"`
	CredentialType string
	FullName       string
	Phone          string
	Email          string
	IDNumber       string
	CreatedAt      string
}
