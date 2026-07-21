package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("api")
	{
		v1 := api.Group("v1")
		{
			user := v1.Group("user")
			{	
				user.POST("/signup", func(ctx *gin.Context) {

				})
			}

			website := v1.Group("website")
			{
				website.POST("/create-website", func(ctx *gin.Context) {
					
				})
			}
		}
	}
}