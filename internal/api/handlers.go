package api

import (
	"net/http"

	"server/internal/models"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

// basic routing for the Rest API
func RegisterRoutes(r *gin.Engine, p *services.PlayerService, b *services.BaseService, bs *services.BuildingService, pv *services.PVPService) {
	r.POST("/guest/login", func(c *gin.Context) {
		resp, err := p.GuestLogin()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	})

	r.GET("/profile", func(c *gin.Context) {
		id := c.MustGet("playerID").(int64)
		resp, err := p.GetProfile(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	})

	r.GET("/base", func(c *gin.Context) {
		id := c.MustGet("playerID").(int64)
		resp, err := b.LoadBase(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resp)
	})

	r.POST("/building/upgrade", func(c *gin.Context) {
		var req models.UpgradeRequest
		if c.BindJSON(&req) != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		id := c.MustGet("playerID").(int64)
		err := bs.StartUpgrade(id, req.BuildingID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	r.POST("/pvp/submit", func(c *gin.Context) {
		var req models.BattleRequest
		if c.BindJSON(&req) != nil {
			c.Status(400)
			return
		}
		id := c.MustGet("playerID").(int64)
		err := pv.SubmitBattle(id, req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"validated": true})
	})
}
