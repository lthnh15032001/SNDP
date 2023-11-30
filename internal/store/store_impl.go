package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lthnh15032001/ngrok-impl/internal/models"
	"github.com/lthnh15032001/ngrok-impl/internal/os"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Storage struct {
	db *gorm.DB
}

// CheckUserExist implements Interface.
func (c *Storage) CheckUserExist(username string, userId string) (*models.UserModel, error) {
	var r models.UserModel
	query := c.db.Model(&models.UserModel{}).Where("username", username).Where("user_id = ?", userId)
	if err := query.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// GetUserACL implements Interface.
func (c *Storage) GetUserACL(userId string, id string) (*models.UserModel, error) {
	var r models.UserModel
	query := c.db.Model(&models.UserModel{}).Where("user_id = ?", userId).Where("id", id)
	if err := query.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// DeleteUser implements Interface.
func (c *Storage) DeleteUserACL(id string) error {
	r := models.UserModel{}
	query := c.db.Model(&models.UserModel{}).Where("id = ?", id)
	if err := query.Delete(&r).Error; err != nil {
		return err
	}
	return nil
}

// GetUser implements Interface.
func (c *Storage) GetAllUsersACL(userId string) (*[]models.UserModel, error) {
	var r []models.UserModel
	query := c.db.Model(&models.UserModel{}).Where("user_id = ?", userId)
	if err := query.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// AddUser implements Interface.
func (c *Storage) AddUserACL(user models.UserModel) error {

	var checkUser models.UserModel

	result := c.db.Where("username = ?", user.Username).Where("user_id = ?", user.UserId).First(&checkUser)
	if result.RowsAffected != 0 {
		// The UserModel does not exist, you might want to handle this case
		return errors.New("UserModel existed")
	}
	if err := c.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// AddUser implements Interface.
func (c *Storage) EditUserACL(id string, user models.UserModel) error {

	var checkUser models.UserModel

	result := c.db.Where("ID = ?", id).Where("user_id = ?", user.UserId).First(&checkUser)
	if result.RowsAffected == 0 {
		// The UserModel does not exist, you might want to handle this case
		return errors.New("User not existed")
	}
	if err := c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "ID"}},
		DoUpdates: clause.AssignmentColumns([]string{"username", "password", "user_remote_policy"}),
	}).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTunnel implements Interface.
func (c *Storage) DeleteTunnel(uuid string) error {
	result := c.db.Where("id = ?", uuid).Delete(&models.TunnelAgentModel{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// GetTunnelActive implements Interface.
func (c *Storage) GetTunnelActive() (*[]models.TunnelAgentModel, error) {
	var r []models.TunnelAgentModel
	query := c.db.Model(&models.TunnelAgentModel{}).Where("status = ?", 1)
	if err := query.Find(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// ChangeTunnelStatus implements Interface.
func (c *Storage) ChangeTunnelStatus(uuid string) error {
	result := c.db.Model(&models.TunnelAgentModel{}).Where("id=?", uuid).Update("status", 0)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// AddTunnel implements Interface.
func (c *Storage) AddTunnel(tunnel models.TunnelAgentModel) error {
	if err := c.db.Create(&tunnel).Error; err != nil {
		return err
	}
	return nil
}

var modelStorage *Storage
var muModelStorage sync.Mutex

func GetOnce() (Interface, *gorm.DB, error) {
	muModelStorage.Lock()
	defer func() {
		muModelStorage.Unlock()
	}()
	if modelStorage != nil {
		return modelStorage, nil, nil
	}
	dbType := os.GetEnv("ENV_DB_TYPE", "mysql")
	host := os.GetEnv("ENV_DB_HOST", "localhost")
	user := os.GetEnv("ENV_DB_USER", "root")
	password := os.GetEnv("ENV_DB_PASSWORD", "my-secret-pw")
	dbName := os.GetEnv("ENV_DB_NAME", "iot")
	switch dbType {
	case "mysql":
		db, err := NewMySQLDB(host, user, password, dbName)

		if err != nil {
			fmt.Printf("err mysql %v\n", err)
			return nil, db, err
		}
		store := New(db)
		return store, db, nil
	default:
		db, err := NewSqliteDB()
		if err != nil {
			fmt.Printf("err sqlite %v/n", err)
			return nil, nil, err
		}
		store := New(db)
		return store, db, nil
	}
}

func New(db *gorm.DB) Interface {
	return &Storage{
		db,
	}
}
