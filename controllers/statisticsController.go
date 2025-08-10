package controllers

import (
	"finance-backend/database"
	"finance-backend/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	userID := c.GetUint("user_id")

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

	// Créer la réponse
	response := gin.H{
		"expenses":   expenses,
		"incomes":    incomes,
		"balance":    balance,
		"categories": categories,
	}

	// Afficher la réponse dans la console
	log.Printf("=== STATISTIQUES POUR L'UTILISATEUR %d ===", userID)
	log.Printf("Dépenses: %.2f", expenses)
	log.Printf("Revenus: %.2f", incomes)
	log.Printf("Solde: %.2f", balance)
	log.Printf("Nombre de catégories: %d", len(categories))

	if len(categories) > 0 {
		log.Printf("Répartition par catégorie:")
		for _, cat := range categories {
			log.Printf("  - %s: %.2f (%.1f%%)",
				cat["category"],
				cat["amount"],
				cat["percent"])
		}
	}
	log.Printf("================================")

	c.JSON(http.StatusOK, response)
}
