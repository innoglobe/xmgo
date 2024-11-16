package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CompanyType represents the type of a company
type CompanyType string

// Company types
const (
	Corporation        CompanyType = "Corporation"
	NonProfit          CompanyType = "NonProfit"
	Cooperative        CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

// Company represents a company entity
type Company struct {
	ID                uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name              string      `gorm:"type:varchar(255);not null" binding:"required" json:"name"`
	Description       string      `gorm:"type:text" json:"description"`
	AmountOfEmployees int         `gorm:"type:int" binding:"required" json:"amount_of_employees"`
	Registered        bool        `gorm:"type:boolean" binding:"required" json:"registered"`
	Type              CompanyType `gorm:"type:varchar(50);not null" binding:"required" json:"type"`
	CreatedAt         time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

// IsValid validates the company type
func (ct CompanyType) IsValid() error {
	switch ct {
	case Corporation, NonProfit, Cooperative, SoleProprietorship:
		return nil
	}
	return fmt.Errorf("invalid company type: %s", ct)
}
