package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lthnh15032001/ngrok-impl/internal/store"
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
	fmt.Println(h.StoreInterface.ListSchedule())
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *UserController) AddUser(c *gin.Context) {
	h.StoreInterface.ListSchedule()
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
