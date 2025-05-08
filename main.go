package main

import (
    "finance-backend/database"
    "finance-backend/routes"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    r.Use(cors.Default()) 
    database.InitDB()
    routes.SetupRoutes(r)
    r.Run(":8080")
}
