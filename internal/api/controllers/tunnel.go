package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lthnh15032001/ngrok-impl/internal/models"
	"github.com/lthnh15032001/ngrok-impl/internal/store"
)

type TunnelController struct {
	StoreInterface store.Interface
}

type AddTunnelDTO struct {
	*models.TunnelAgentModel
}

func (h *TunnelController) AddTunnel(c *gin.Context) {
	m := AddTunnelDTO{}
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}

	sInterface := h.StoreInterface
	uuidV4, _ := uuid.NewRandom()
	s := models.TunnelAgentModel{
		ID:           uuidV4.String(),
		Name:         m.Name,
		Region:       m.Region,
		IP:           m.IP,
		Version:      m.Version,
		TunnelOnline: m.TunnelOnline,
		StartedAt:    time.Now(),
		OS:           m.OS,
		Metadata:     m.Metadata,
	}
	err := sInterface.AddTunnel(s)
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

func (h *TunnelController) GetTunnelActive(c *gin.Context) {
	allTunnel, err := h.StoreInterface.GetTunnelActive()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   allTunnel,
	})
}
