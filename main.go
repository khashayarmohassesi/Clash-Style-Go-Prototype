package main

import (
	"log"

	"server/db"
	"server/internal/api"
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.NewPostgres()
	defer database.Close()

	playerSvc := services.NewPlayerService(database)
	baseSvc := services.NewBaseService(database)
	buildingSvc := services.NewBuildingService(database)
	pvpSvc := services.NewPVPService(database)

	r := gin.Default()
	r.Use(api.PlayerMiddleware(database))
	api.RegisterRoutes(r, playerSvc, baseSvc, buildingSvc, pvpSvc)

	log.Fatal(r.Run(":8080"))
}
