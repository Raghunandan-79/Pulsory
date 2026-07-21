package handlers

import (
	"net/http"
	"time"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/Raghunandan-79/pulsory/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateWebsiteRequest struct {
	URL string `json:"url" binding:"required"`
}

func CreateWebsite(ctx *gin.Context) {
	var req CreateWebsiteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	website := models.Website{
		URL: req.URL,
		TimeAdded: time.Now(),
		UserID: uuid.MustParse(userID.(string)),
	}

	if err := config.DB.Create(&website).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create website",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"website": website,
	})
}

func GetWebsiteById(ctx *gin.Context) {

}