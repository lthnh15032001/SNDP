package models

import (
	"time"

	"gorm.io/datatypes"
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
	UserRemoteId string    `json:"userRemoteId" gorm:"size:191"`
}

type UserRemotePolicy struct {
}
type UserModel struct {
	ID               string         `gorm:"primaryKey"`
	UserId           string         `json:"userId" `  // keycloak userid
	Username         string         `json:"username"` // composite primary key
	Password         string         `json:"password"`
	UserRemotePolicy datatypes.JSON `json:"userRemotePolicy" gorm:"type:json"`
}
