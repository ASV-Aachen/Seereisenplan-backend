package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type License struct {
	Name        string
	Description string
	ShortName   string
}

type LicenseShare struct {
	Type          License
	Since         time.Time
	LicenseNumber string
}

type Sailor struct {
	First_name string
	Last_name  string
	ID         string `gorm:"->;<-:create;primaryKey"`
	License    []LicenseShare
	Guest      bool
}

func (sailor *Sailor) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	sailor.ID = uuid.NewString()
	return
}
