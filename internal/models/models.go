package models

import (
	"time"

	_ "gorm.io/datatypes"
)

type TunnelAgentModel struct {
	ID           string    `gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:191"`
	Region       string    `json:"region"`
	IP           string    `json:"ip"`
	Version      string    `json:"version"`
	TunnelOnline int32     `json:"tunnelonline"`
	StartedAt    time.Time `json:"startedAt", gorm:"column:startedAt"`
	OS           string    `json:"os"`
	Metadata     string    `json:"metadata"`
	Status       int       `json:"status" gorm:"default:1"`
}
