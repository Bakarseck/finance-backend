package controllers

import (
    "finance-backend/database"
    "finance-backend/models"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("verysecretkey")

func RegisterUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Vérifier si email existe déjà
    var existing models.User
    if err := database.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email déjà utilisé"+err.Error()})
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
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"+err.Error()})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"+err.Error()})
        return
    }

    // ✅ Création du token JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
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
