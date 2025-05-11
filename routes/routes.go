package routes

import (
	"finance-backend/controllers"
	"finance-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Routes publiques (non protégées)
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)

		// Routes protégées
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/transactions", controllers.GetTransactions)
			protected.POST("/transactions", controllers.CreateTransaction)
			protected.DELETE("/transactions/:id", controllers.DeleteTransaction)

			protected.GET("/categories", controllers.GetCategories)
			protected.POST("/categories", controllers.CreateCategory)

			protected.GET("/balance", controllers.GetBalance)

			protected.GET("/balance", controllers.GetBalance)
			protected.GET("/statistics", controllers.GetStatistics)
		}
	}
}
