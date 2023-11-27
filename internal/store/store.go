package store

import (
	"github.com/lthnh15032001/ngrok-impl/internal/models"
)

//go:generate mockery --name Interface
type Interface interface {
	AddTunnel(tunnel models.TunnelAgentModel) error
	ChangeTunnelStatus(uuid string) error
	GetTunnelActive() (*[]models.TunnelAgentModel, error)
	DeleteTunnel(uuid string) error

	AddUser(user models.UserModel) error
}
