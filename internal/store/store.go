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

	AddUserACL(user models.UserModel) error
	GetAllUsersACL(userid string) (*[]models.UserModel, error)
	GetUserACL(userid string, id string) (*models.UserModel, error)
	DeleteUserACL(userid string) error
	EditUserACL(id string, user models.UserModel) error
}
