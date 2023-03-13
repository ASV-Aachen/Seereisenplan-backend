package db

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Cruise struct {
	ID                string `gorm:"->;<-:create;primaryKey"`
	CruiseName        string
	CuriseDescription string
	StartDate         time.Time
	EndDate           time.Time
	StartPort         string
	EndPort           string
	MaxBerths         int `gorm:"default:0"`
	Season            int
}

func (cruise *Cruise) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	season := cruise.StartDate.Year()
	cruise.ID = uuid.NewString()
	cruise.Season = season
	return
}

type CruiseShare struct {
	Cruise   Cruise `gorm:"constraint:OnDelete:CASCADE"`
	Sailor   Sailor
	Position string  `gorm:"default:CREW"`
	Distance float64 `gorm:"default:0"`
}
