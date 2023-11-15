package store

import (
	"github.com/lthnh15032001/ngrok-impl/internal/models"
)

//go:generate mockery --name Interface
type Interface interface {
	SavePolicy(policy models.Policy) error
	GetPolicyByName(name string) (*models.Policy, error)
	GetPolicyBySchedule(name string) (*[]models.Policy, error)
	ListPolicy() (*[]models.Policy, error)
	ListPolicyByProvider(name string) (*[]models.Policy, error)
	DeletePolicy(name string) error
	SaveSchedule(schedule models.ScheduleModel) error
	GetSchedule(name string) (*models.ScheduleModel, error)
	ListSchedule() (*[]models.ScheduleModel, error)
	DeleteSchedule(name string) error
}
