package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	Model
	LongURL      string    `gorm:"not null" json:"long_url"`
	ShortURL     string    `gorm:"unique;not null" json:"short_url"`
	LastAccessed time.Time `json:"last_accessed"`
	AccessCount  int       `gorm:"default:0" json:"access_count"`
}

func (a *URL) CreateURL(db *gorm.DB) error {
	err := db.
		Create(a).
		Error

	if err != nil {
		return err
	}

	return nil
}

func GetURL(db *gorm.DB, shortCode string) (*URL, error) {
	var url URL
	err := db.
		Where("short_url = ?", shortCode).
		First(&url).
		Error

	if err != nil {
		return nil, err
	}

	return &url, nil
}
