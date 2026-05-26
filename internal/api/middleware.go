package api

//Using Gin package for implementation of Rest APIs
import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PlayerMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/guest/login" {
			c.Next()
			return
		}
		v := c.GetHeader("X-Player-ID")
		if v == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing X-Player-ID"})
			return
		}
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		var exists bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM players WHERE id=$1)`, id).Scan(&exists)
		if err != nil || !exists {
			c.AbortWithStatus(401)
			return
		}
		c.Set("playerID", id)
		c.Next()
	}
}
