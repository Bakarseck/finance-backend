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

	log.Printf("=== DEBUG STATISTIQUES ===")
	log.Printf("Recherche pour l'utilisateur: %d", userID)
	log.Printf("Période: du %s au %s", monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02"))

	var transactions []models.Transaction
	result := database.DB.Where("user_id = ? AND date >= ? AND date < ?", userID, monthStart, monthEnd).Find(&transactions)

	log.Printf("Requête SQL exécutée")
	log.Printf("Nombre de transactions trouvées: %d", len(transactions))
	log.Printf("Erreur de base de données: %v", result.Error)

	// Afficher les transactions trouvées pour debug
	if len(transactions) > 0 {
		log.Printf("Transactions trouvées:")
		for i, t := range transactions {
			log.Printf("  [%d] ID: %d, Type: %s, Montant: %.2f, Catégorie: %s, Date: %s",
				i+1, t.ID, t.Type, t.Amount, t.Category, t.Date.Format("2006-01-02"))
		}
	} else {
		log.Printf("AUCUNE TRANSACTION TROUVÉE pour cette période!")

		// Vérifier s'il y a des transactions pour cet utilisateur (toutes périodes)
		var allTransactions []models.Transaction
		database.DB.Where("user_id = ?", userID).Find(&allTransactions)
		log.Printf("Total transactions pour cet utilisateur (toutes périodes): %d", len(allTransactions))

		if len(allTransactions) > 0 {
			log.Printf("Exemples de transactions existantes:")
			for i, t := range allTransactions[:5] { // Afficher les 5 premières
				log.Printf("  [%d] ID: %d, Type: %s, Montant: %.2f, Catégorie: %s, Date: %s",
					i+1, t.ID, t.Type, t.Amount, t.Category, t.Date.Format("2006-01-02"))
			}
		}
	}

	expenses := 0.0
	incomes := 0.0
	balance := 0.0
	categoryTotals := make(map[string]float64)
	categoryLabels := make(map[string]string)

	for _, t := range transactions {
		log.Printf("Traitement transaction: Type=%s, Montant=%.2f, Catégorie=%s", t.Type, t.Amount, t.Category)

		if t.Type == "expense" {
			expenses += t.Amount
			categoryTotals[t.Category] += t.Amount
			categoryLabels[t.Category] = t.Category
			log.Printf("  -> Ajouté aux dépenses: %.2f", t.Amount)
		} else if t.Type == "income" {
			incomes += t.Amount
			log.Printf("  -> Ajouté aux revenus: %.2f", t.Amount)
		} else {
			log.Printf("  -> Type inconnu: %s", t.Type)
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
