/*
Author: Jason Whatson
Date: 23/10/2022
*/
package product

import (
	"database/sql"
)

var dbConn *sql.DB

type Product struct {
	SKU           string  `json:"sku"`
	Barcode       string  `json:"barcode"`
	Name          string  `json:"name"`
	StandardPrice float32 `json:"standardprice" binding:"required"`
	SpecialPrice  float32 `json:"specialprice"`
}

type ProductPriceUpdateStatus struct {
	Barcode string `json:"barcode"`
	Status  string `json:"status"`
}

// Returns all products in our database
func GetProducts() []Product {

	// Get all Products which have a valid Standard Price
	rows, err := dbConn.Query("SELECT * FROM Products WHERE StandardPrice is not null AND StandardPrice != ''")
	checkErr(err)

	defer rows.Close() //good habit to close

	products := []Product{}

	var specialPrice sql.NullFloat64

	for rows.Next() {
		prod := Product{}
		err = rows.Scan(&prod.SKU, &prod.Barcode, &prod.Name, &prod.StandardPrice, &specialPrice)
		if specialPrice.Valid { //Handle null / empty special price in DB
			prod.SpecialPrice = float32(specialPrice.Float64)
		} else {
			prod.SpecialPrice = 0
		}
		checkErr(err)
		products = append(products, prod)
	}

	return products

}

// Updates products pricing if needed and returns an a array of status
func UpdateProductsPricing(products []Product) []ProductPriceUpdateStatus {

	var productPriceUpdateStatus = []ProductPriceUpdateStatus{}

	var status string

	for _, product := range products {

		priceNeedsUpdating := doesRecordNeedPriceUpdate(product.Barcode, product.StandardPrice, product.SpecialPrice)

		if priceNeedsUpdating {

			stmt, err := dbConn.Prepare("UPDATE Products SET StandardPrice=?, SpecialPrice=? WHERE barcode=?")
			checkErr(err)

			res, errExec := stmt.Exec(product.StandardPrice, product.SpecialPrice, product.Barcode)
			checkErr(errExec)

			affect, err := res.RowsAffected()
			checkErr(err)

			if affect >= 1 {
				status = "updated"
			} else {
				status = "notfound" //Barcode not found
			}

		} else {
			status = "ignored"
		}

		productPriceUpdateStatus = append(productPriceUpdateStatus, ProductPriceUpdateStatus{
			Barcode: product.Barcode,
			Status:  status,
		})

	}

	return productPriceUpdateStatus

}

// Checks if the products record needs a price update
// Returns true if count <= 0 (found same record). Otherwise returns false
func doesRecordNeedPriceUpdate(barcode string, StandardPrice float32, SpecialPrice float32) bool {

	count := -1

	rows, err := dbConn.Query(
		"SELECT count(*) count FROM Products WHERE StandardPrice = $1 AND SpecialPrice = $2 AND Barcode=? ",
		StandardPrice, SpecialPrice, barcode,
	)

	checkErr(err)

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&count)
	checkErr(err)

	return (count <= 0)

}

// Connects to sqlite3 database
func ConnectToDatabase() {
	var err error
	dbConn, err = sql.Open("sqlite3", "../products.db")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
