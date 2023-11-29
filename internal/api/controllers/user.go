package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lthnh15032001/ngrok-impl/internal/models"
	"github.com/lthnh15032001/ngrok-impl/internal/store"
	"gorm.io/datatypes"
)

type UserController struct {
	StoreInterface store.Interface
}

// Status godoc
// @Summary Responds with 200 if service is running
// @Description health check for service
// @Produce  json
// @Success 200 {string} pong
// @Router /health/ping [get]
func (h *UserController) GetAllUsers(c *gin.Context) {
	getUsers, err := h.StoreInterface.GetAllUsersACL(c.GetString("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   getUsers,
	})
}

func (h *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	getUser, err := h.StoreInterface.GetUserACL(c.GetString("userid"), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   getUser,
	})
}

type AddUserDTO struct {
	*models.UserModel
}

func (h *UserController) AddUserAuthen(c *gin.Context) {
	m := AddUserDTO{}

	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	uuidV4, _ := uuid.NewRandom()

	sInterface := h.StoreInterface
	s := models.UserModel{
		ID:               uuidV4.String(),
		UserId:           c.GetString("userid"),
		Username:         m.Username,
		Password:         m.Password,
		UserRemotePolicy: datatypes.JSON(m.UserRemotePolicy),
	}
	err := sInterface.AddUserACL(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   s,
	})
}

func (h *UserController) EditUserAuthen(c *gin.Context) {
	m := AddUserDTO{}
	id := c.Param("id")
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}

	sInterface := h.StoreInterface
	s := models.UserModel{
		ID:               id,
		UserId:           c.GetString("userid"),
		Username:         m.Username,
		Password:         m.Password,
		UserRemotePolicy: datatypes.JSON(m.UserRemotePolicy),
	}
	err := sInterface.EditUserACL(id, s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   s,
	})
}

func (h *UserController) DeleteUserAuthen(c *gin.Context) {
	sInterface := h.StoreInterface
	id, ok := c.GetQuery("id")
	if id == "" && ok == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   "Missing id",
		})
		return
	}
	err := sInterface.DeleteUserACL(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   " s",
	})
	return
}
