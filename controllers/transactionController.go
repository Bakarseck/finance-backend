package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"finance-backend/database"
	"finance-backend/models"

	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	userID := c.GetUint("user_id")
	database.DB.Where("user_id = ?", userID).Find(&transactions)
	c.JSON(http.StatusOK, transactions)
}

func CreateTransaction(c *gin.Context) {
	// Log le JSON brut reçu
	body, _ := c.GetRawData()
	fmt.Printf(">> JSON reçu brut : %s\n", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Struct temporaire pour désérialiser correctement
	var req struct {
		Icon       string `json:"icon"`
		Title      string `json:"title"`
		Subtitle   string `json:"subtitle"`
		Date       string `json:"date"`
		Amount     string `json:"amount"`
		Type       string `json:"type"`
		Category   string `json:"category"`
		CategoryID uint   `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Erreur de liaison JSON : %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		fmt.Printf("Erreur parsing date: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date invalide. Format attendu : YYYY-MM-DD"})
		return
	}

	// Nettoyage de l'amount
	var amountValue float64
	fmt.Sscanf(req.Amount, "%f", &amountValue)

	// Récupérer l'ID de l'utilisateur depuis le contexte
	userID := c.GetUint("user_id")

	transaction := models.Transaction{
		Icon:       req.Icon,
		Title:      req.Title,
		Subtitle:   req.Subtitle,
		Date:       parsedDate,
		Amount:     amountValue,
		Type:       req.Type,
		Category:   req.Category,
		CategoryID: req.CategoryID,
		UserID:     userID,
	}

	database.DB.Create(&transaction)
	c.JSON(http.StatusCreated, transaction)
}

func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{})
	c.Status(http.StatusNoContent)
}

func GetBalance(c *gin.Context) {
	var transactions []models.Transaction
	userID := c.GetUint("user_id")
	database.DB.Where("user_id = ?", userID).Find(&transactions)

	var balance float64 = 0
	for _, t := range transactions {
		balance += t.Amount
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
