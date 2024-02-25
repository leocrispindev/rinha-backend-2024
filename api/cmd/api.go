package main

import (
	"rinha-backend-2024/api/internal/controller"
	"rinha-backend-2024/api/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	r := gin.Default()
	r.POST("/clientes/:id/transacoes", controller.HandlerTransaction)
	r.GET("/clientes/:id/extrato", controller.HandlerExtract)
	r.Run(":9999") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
