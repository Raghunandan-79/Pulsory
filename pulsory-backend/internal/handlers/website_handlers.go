package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Raghunandan-79/pulsory/internal/config"
	"github.com/Raghunandan-79/pulsory/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	url := strings.TrimSpace(req.URL)
	parsedUserID := uuid.MustParse(userID.(string))

	var existing models.Website
	err := config.DB.Where("url = ? AND user_id = ?", url, parsedUserID).First(&existing).Error
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Website already exists",
		})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check existing website",
		})
		return
	}

	website := models.Website{
		URL:       url,
		TimeAdded: time.Now(),
		UserID:    parsedUserID,
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
	websiteID := strings.TrimSpace(ctx.Param("websiteId"))

	if _, err := uuid.Parse(websiteID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid website id",
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

	var website models.Website
	err := config.DB.Where("id = ? AND user_id = ?", websiteID, userID).First(&website).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Website not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch website",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"website": website,
	})
}