package database

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "finance-backend/models"
)

var DB *gorm.DB

func InitDB() {
    var err error
    DB, err = gorm.Open(sqlite.Open("finance.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Échec connexion à la base : ", err)
    }

    DB.AutoMigrate(&models.Transaction{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
}
