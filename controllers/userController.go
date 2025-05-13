package controllers

import (
	"finance-backend/database"
	"finance-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte("verysecretkey")

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si email existe déjà
	var existing models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email déjà utilisé"})
		return
	}

	// Hacher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de sécurité"})
		return
	}

	user.Password = string(hashedPassword)

	database.DB.Create(&user)
	c.JSON(http.StatusCreated, gin.H{
		"id":        user.ID,
		"full_name": user.FullName,
		"email":     user.Email,
		"currency":  user.Currency,
	})
}

func LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect" + err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect" + err.Error()})
		return
	}

	// ✅ Création du token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"full_name": user.FullName,
		"currency":  user.Currency,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"email":     user.Email,
			"currency":  user.Currency,
		},
	})
}

func ProtectedRoute(c *gin.Context) {
	userID := c.GetUint("user_id")
	// Use userID as needed
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"full_name": user.FullName,
		"email":     user.Email,
		"currency":  user.Currency,
	})
}

func UpdateName(c *gin.Context) {
	userID := c.GetUint("user_id")

	var payload struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nom invalide"})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("full_name", payload.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de mise à jour du nom"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nom mis à jour avec succès"})
}

func UpdateEmail(c *gin.Context) {
	userID := c.GetUint("user_id")

	var payload struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email invalide"})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("email", payload.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de mise à jour de l'email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email mis à jour avec succès"})
}

func UpdateCurrency(c *gin.Context) {
	userID := c.GetUint("user_id")

	var payload struct {
		Currency string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Devise invalide"})
		return
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("currency", payload.Currency).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de mise à jour de la devise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Devise mise à jour avec succès"})
}

