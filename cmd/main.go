package main

import (
	"products-api/controller"
	"products-api/db"
	"products-api/repository"
	"products-api/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	dbConnection, err := db.ConnectDB()

	if err != nil {
		panic(err)
	}

	ProductRepository := repository.NewProductRepository(dbConnection)
	ProductUsecase := usecase.NewProductUsecase(ProductRepository)
	ProductController := controller.NewProductContrtoller(ProductUsecase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/products", ProductController.GetProducts)
	server.GET("products/:productId", ProductController.GetProduct)
	server.POST("/products", ProductController.CreateProduct)
	server.PUT("/products/:id", ProductController.UpdateProduct)

	server.Run(":8000")
}
