package main

import (
	"log"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/Raghunandan-79/pulsory/internal/models"
	"github.com/Raghunandan-79/pulsory/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Region{},
		&models.Website{},
		&models.WebsiteTick{},
	)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default();
	routes.RegisterRoutes(router)
	
	router.Run(":8080")
}