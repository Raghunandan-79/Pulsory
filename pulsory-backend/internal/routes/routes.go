package routes

import (
	"github.com/Raghunandan-79/pulsory/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			user := v1.Group("user")
			{	
				user.POST("/signup", handlers.Signup)
				user.POST("/signin", handlers.Signin)
			}

			website := v1.Group("website")
			{
				website.POST("/create-website", func(ctx *gin.Context) {

				})
			}
		}
	}
}