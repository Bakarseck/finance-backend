package routes

import (
    "finance-backend/controllers"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        api.GET("/transactions", controllers.GetTransactions)
        api.POST("/transactions", controllers.CreateTransaction)
        api.DELETE("/transactions/:id", controllers.DeleteTransaction)

		api.GET("/categories", controllers.GetCategories)
		api.POST("/categories", controllers.CreateCategory)

        api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)
        api.GET("/balance", controllers.GetBalance)
    }
}
