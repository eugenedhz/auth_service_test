package models

type Session struct {
	UserID           string `gorm:"primary_key"` // GUID
	RefreshTokenHash []byte
}
