package services

import (
	"fmt"
	"github.com/VaheMuradyan/Sport/models"
	"gorm.io/gorm"
	"time"
)

type CoefficientService struct {
	db *gorm.DB
}

func NewCoefficientService(db *gorm.DB) *CoefficientService {
	return &CoefficientService{
		db: db,
	}
}

func (s *CoefficientService) UpdateMarketCoefficient(marketID uint, newCoefficient float64, userID uint) (*models.CoefficientUpdateResponse, error) {
	var market models.Market
	if err := s.db.First(&market, marketID).Error; err != nil {
		return nil, fmt.Errorf("market not found: %v", err)
	}

	if newCoefficient < market.MinCoefficient || newCoefficient > market.MaxCoefficient {
		return nil, fmt.Errorf("odds %.2f out of bounds [%.2f, %.2f]", newCoefficient, market.MinCoefficient, market.MaxCoefficient)
	}

	oldCoefficient := market.CurrentCoefficient

	tx := s.db.Begin()

	market.PreviousCoefficient = oldCoefficient
	market.CurrentCoefficient = newCoefficient
	market.LastUpdated = time.Now()

	if err := tx.Save(&market).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update market: %v", err)
	}

	coefficientHistory := models.CoefficientHistory{
		MarketID:    marketID,
		OldValue:    oldCoefficient,
		NewValue:    newCoefficient,
		ChangedByID: userID,
		Timestamp:   time.Now(),
	}

	if err := tx.Create(&coefficientHistory).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to record coefficient history: %v", err)
	}

	tx.Commit()

	return &models.CoefficientUpdateResponse{
		Success:        true,
		Message:        "Coefficient updated successfully",
		MarketID:       marketID,
		OldCoefficient: oldCoefficient,
		NewCoefficient: newCoefficient,
		UpdatedAt:      market.LastUpdated,
	}, nil
}

func (s *CoefficientService) GetMarketWithHistory(marketID uint) (*models.Market, error) {
	var market models.Market
	err := s.db.Preload("CoefficientHistory").Preload("Event").First(&market, marketID).Error
	return &market, err
}
