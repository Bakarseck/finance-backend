package controllers

import (
	"finance-backend/database"
	"finance-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer l'ID de l'utilisateur depuis le contexte
	userID := c.GetUint("user_id")
	category.UserID = userID

	database.DB.Create(&category)
	c.JSON(http.StatusCreated, category)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	userID := c.GetUint("user_id")
	database.DB.Where("user_id = ?", userID).Find(&categories)
	c.JSON(http.StatusOK, categories)
}
