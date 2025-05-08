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

    database.DB.Create(&category)
    c.JSON(http.StatusCreated, category)
}

func GetCategories(c *gin.Context) {
    var categories []models.Category
    database.DB.Find(&categories)
    c.JSON(http.StatusOK, categories)
}
