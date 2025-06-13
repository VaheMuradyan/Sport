package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
}

type Country struct {
	gorm.Model
	Name         string
	Competitions []Competition `gorm:"foreignKey:CountryID" json:"competitions,omitempty"`
	Code         string        `gorm:"unique;size:3" json:"code"`
}

type Competition struct {
	gorm.Model
	Name      string  `json:"name"`
	CountryID uint    `json:"country_id"`
	Country   Country `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	SportID   uint    `json:"sport_id"`
	Sport     Sport   `gorm:"foreignKey:SportID" json:"sport,omitempty"`
	Teams     []Team  `gorm:"many2many:competition_teams;" json:"teams,omitempty"`
	Events    []Event `gorm:"foreignKey:CompetitionID" json:"events,omitempty"`
	Active    bool    `gorm:"default:true" json:"active"`
}

type Market struct {
	gorm.Model
	Name                string               `json:"name"`
	Type                string               `json:"type"` // match_winner, over_under, handicap, etc.
	EventID             uint                 `json:"event_id"`
	Event               Event                `gorm:"foreignKey:EventID" json:"event,omitempty"`
	CurrentCoefficient  float64              `gorm:"type:decimal(9,4);" json:"current_coefficient"`
	PreviousCoefficient float64              `gorm:"type:decimal(9,4);" json:"previous_coefficient"`
	MinCoefficient      float64              `gorm:"type:decimal(9,4);default:1.01" json:"min_coefficient"`
	MaxCoefficient      float64              `gorm:"type:decimal(9,4);default:100.00" json:"max_coefficient"`
	Active              bool                 `gorm:"default:true" json:"active"`
	CoefficientHistory  []CoefficientHistory `gorm:"foreignKey:MarketID" json:"coefficient_history,omitempty"`
	LastUpdated         time.Time            `json:"last_updated"`
}

type CoefficientHistory struct {
	gorm.Model
	MarketID    uint      `json:"market_id"`
	Market      Market    `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	OldValue    float64   `gorm:"type:decimal(9,4)" json:"old_value"`
	NewValue    float64   `gorm:"type:decimal(9,4)" json:"new_value"`
	ChangedByID uint      `json:"changed_by_id"`
	ChangedBy   User      `gorm:"foreignKey:ChangedByID" json:"changed_by,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

type Event struct {
	gorm.Model
	Name          string      `json:"name"`
	CompetitionID uint        `json:"competition_id"`
	Competition   Competition `gorm:"foreignKey:CompetitionID" json:"competition,omitempty"`
	Markets       []Market    `gorm:"foreignKey:EventID" json:"markets,omitempty"`
	Teams         []Team      `gorm:"many2many:event_teams;" json:"teams,omitempty"`
	StartTime     time.Time   `json:"start_time"`
	Status        string      `gorm:"default:'scheduled'" json:"status"`
	IsLive        bool        `gorm:"default:false" json:"is_live"`
}

type Team struct {
	gorm.Model
	Name         string        `json:"name"`
	Rating       int           `json:"rating"`
	CountryID    uint          `json:"country_id"`
	Country      Country       `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	Competitions []Competition `gorm:"many2many:competition_teams;" json:"competitions,omitempty"`
	Events       []Event       `gorm:"many2many:event_teams;" json:"events,omitempty"`
	Sports       []Sport       `gorm:"many2many:sport_teams;" json:"sports,omitempty"`
}

type Sport struct {
	gorm.Model
	Name         string        `json:"name"`
	Competitions []Competition `gorm:"foreignKey:SportID" json:"competitions,omitempty"`
	Teams        []Team        `gorm:"many2many:sport_teams;" json:"teams,omitempty"`
	Code         string        `gorm:"unique" json:"code"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CoefficientUpdateRequest struct {
	MarketID       uint    `json:"market_id" binding:"required"`
	NewCoefficient float64 `json:"new_Coefficient" binding:"required"`
}

type CoefficientUpdateResponse struct {
	Success        bool      `json:"success"`
	Message        string    `json:"message"`
	MarketID       uint      `json:"market_id"`
	OldCoefficient float64   `json:"old_coefficient"`
	NewCoefficient float64   `json:"new_coefficient"`
	UpdatedAt      time.Time `json:"updated_at"`
}
