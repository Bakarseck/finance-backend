package controllers

import (
	"finance-backend/database"
	"finance-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Début et fin du mois courant
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	var transactions []models.Transaction
	database.DB.Where("user_id = ? AND date >= ? AND date < ?", userID, monthStart, monthEnd).Find(&transactions)

	expenses := 0.0
	incomes := 0.0
	balance := 0.0
	categoryTotals := make(map[string]float64)
	categoryLabels := make(map[string]string)

	for _, t := range transactions {
		if t.Type == "expense" {
			expenses += t.Amount
			categoryTotals[t.Category] += t.Amount
			categoryLabels[t.Category] = t.Category
		} else if t.Type == "income" {
			incomes += t.Amount
		}
		balance += t.Amount
	}

	// Préparer la répartition par catégorie
	totalExpenses := expenses
	categories := []gin.H{}
	for cat, amount := range categoryTotals {
		percent := 0.0
		if totalExpenses > 0 {
			percent = (amount / totalExpenses) * 100
		}
		categories = append(categories, gin.H{
			"category": cat,
			"amount":   amount,
			"percent":  percent,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"expenses":   expenses,
		"incomes":    incomes,
		"balance":    balance,
		"categories": categories,
	})
}
