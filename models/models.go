package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

type Country struct {
	gorm.Model
	Name         string
	Competitions []Competition `gorm:"foreignKey:CountryID"`
}

type Competition struct {
	gorm.Model
	Name      string
	CountryID uint
	Country   Country `gorm:"foreignKey:CountryID"`
	Teams     []Team  `gorm:"many2many:competition_teams;"`
	Events    []Event `gorm:"foreignKey:CompetitionID"`
}

type Market struct {
	gorm.Model
	Name    string
	EventID uint
	Event   Event `gorm:"foreignKey:EventID"`
}

type Event struct {
	gorm.Model
	Name          string
	CompetitionID uint
	Competition   Competition `gorm:"foreignKey:CompetitionID"`
	Markets       []Market    `gorm:"foreignKey:EventID"`
	Teams         []Team      `gorm:"many2many:event_teams;"`
}

type Team struct {
	gorm.Model
	Name         string
	Rating       bool
	CountryID    uint
	Country      Country       `gorm:"foreignKey:CountryID"`
	Competitions []Competition `gorm:"many2many:competition_teams;"`
	Events       []Event       `gorm:"many2many:event_teams;"`
	Sports       []Sport       `gorm:"many2many:sport_teams;"`
}

type Sport struct {
	gorm.Model
	Name  string
	Teams []Team `gorm:"many2many:sport_teams;"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
