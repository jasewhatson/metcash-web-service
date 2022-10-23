/*
Author: Jason Whatson
Date: 23/10/2022
*/
package main

import (
	"metcashwebservice/src/product"
	"net/http"

	//3rd party imports
	"github.com/gin-gonic/gin"      //Provides REST API
	_ "github.com/mattn/go-sqlite3" //Provide data persistence via Sqlite
)

func main() {
	product.ConnectToDatabase()
	router := gin.Default()
	router.GET("/products", getProducts)
	router.POST("/pricing", postPricing)
	router.Run("localhost:8880")
}

// Responds with the list of all products as an JSON array.
func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, product.GetProducts())
}

// postPricing adds an album from JSON received in the request body.
func postPricing(c *gin.Context) {

	var productPricing []product.Product

	// Call BindJSON to bind the received JSON to Product
	if err := c.BindJSON(&productPricing); err != nil {
		return
	}

	productPriceUpdateStatus := product.UpdateProductsPricing(productPricing)

	c.IndentedJSON(http.StatusOK, productPriceUpdateStatus)
}
