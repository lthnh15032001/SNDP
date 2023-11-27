package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lthnh15032001/ngrok-impl/internal/models"
	"github.com/lthnh15032001/ngrok-impl/internal/os"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

// AddUser implements Interface.
func (c *Storage) AddUser(user models.UserModel) error {

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

// func (c *Storage) SavePolicy(m models.Policy) error {

// 	if err := c.db.Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "name"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"name", "display_name", "projects", "tags", "schedule_name", "provider"}),
// 	}).Create(&m).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Storage) GetPolicyByName(name string) (*models.Policy, error) {
// 	r := models.Policy{}
// 	query := c.db.Model(&models.Policy{}).Where("name=?", name)
// 	if err := query.First(&r).Error; err != nil {
// 		return nil, err
// 	}
// 	return &r, nil
// }

// func (c *Storage) GetPolicyBySchedule(name string) (*[]models.Policy, error) {
// 	var r []models.Policy
// 	query := c.db.Model(&models.Policy{}).Where("schedule_name=?", name)
// 	if err := query.Find(&r).Error; err != nil {
// 		return nil, err
// 	}
// 	return &r, nil
// }

// func (c *Storage) ListPolicy() (*[]models.Policy, error) {
// 	var r []models.Policy
// 	if err := c.db.Find(&r).Error; err != nil {
// 		return nil, err
// 	}
// 	return &r, nil
// }

// func (c *Storage) ListPolicyByProvider(name string) (*[]models.Policy, error) {
// 	var r []models.Policy
// 	if err := c.db.Where("name <> ?", name).Find(&r).Error; err != nil {
// 		return nil, err
// 	}
// 	return &r, nil
// }

// func (c *Storage) DeletePolicy(name string) error {
// 	r := models.Policy{}
// 	query := c.db.Model(&models.Policy{}).Where("name=?", name)
// 	if err := query.Delete(&r).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Storage) SaveSchedule(m models.ScheduleModel) error {
// 	if err := c.db.Clauses(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "name"}},
// 		DoUpdates: clause.AssignmentColumns([]string{"name", "display_name", "time_zone", "schedule"}),
// 	}).Create(&m).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Storage) GetSchedule(name string) (*models.ScheduleModel, error) {
// 	r := models.ScheduleModel{}
// 	query := c.db.Model(&models.ScheduleModel{}).Where("name = ?", name)
// 	if err := query.Find(&r).Error; err != nil {
// 		return nil, err
// 	}
// 	return &r, nil
// }

// func (c *Storage) ListSchedule() (*[]models.ScheduleModel, error) {
// 	var r []models.ScheduleModel
// 	if err := c.db.Find(&r).Error; err != nil {
// 		return nil, err
// 	}
// 	return &r, nil
// }

// func (c *Storage) DeleteSchedule(name string) error {
// 	r := models.ScheduleModel{}
// 	query := c.db.Model(&models.ScheduleModel{}).Where("name=?", name)
// 	if err := query.Delete(&r).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
