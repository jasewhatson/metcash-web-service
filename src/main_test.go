/*
Author: Jason Whatson
Date: 23/10/2022
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"metcashwebservice/src/product"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var updatedPrice string

func init() {
	//For each time we run the test
	//Set a new unique price for 9300675009775 (Coca Cola Diet Coke) in jsonData for TestPostPricing
	//TestGetProducts will check that we are getting this value back from our db via the web service
	updatedPrice = fmt.Sprint(getNextUniquePrice())
}

// Test POST /pricing
func TestPostPricing(t *testing.T) {

	product.ConnectToDatabase()

	// Switch to test mode to turn off verbose logging
	gin.SetMode(gin.TestMode)

	// Setup our router, just like we did in the main function, and
	// register our routes
	r := gin.Default()
	r.POST("/pricing", postPricing)

	var jsonData = []byte(`[
		{
		  "barcode": "9313820004501",
		  "name": "7Up Lemonade Soda, 1.25 Litre",
		  "standardprice": 2.8,
		  "specialprice": 2.5
		},
		{
		  "barcode": "9300675009775",
		  "name": "Coca Cola Diet Coke",
		  "standardprice": ` + updatedPrice + `,
		  "specialprice": null
		},
		{
		  "barcode": "9339423008845",
		  "name": "Abbott's Bakery Sourdough Rye Loaf , 760 Gram",
		  "standardprice": 6.75,
		  "specialprice": 4.50
		},
		{
		  "barcode": "9300619513023",
		  "name": "Aeroplane Jelly Lite Lime, 2 Each",
		  "standardprice": 2.45,
		  "specialprice": 2
		},
		{
		  "barcode": "9300632064168",
		  "name": "Ajax Professional Antibacterial Disinfectant Bathroom Clean, 500 Millilitre",
		  "standardprice": 5.8,
		  "specialprice": null
		}
	   ]`)

	// Create the mock request that we would like to test.
	req, err := http.NewRequest(http.MethodPost, "/pricing", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check for HTTP 200
	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	} else {
		responseData, _ := io.ReadAll(w.Body)

		var updateStatus []product.ProductPriceUpdateStatus
		json.Unmarshal(responseData, &updateStatus)

		for _, v := range updateStatus {

			//Pricing for Coca Cola Diet Coke (barcode 9300675009775) gets updated with a new price every run
			if v.Barcode == "9300675009775" {
				if v.Status != "updated" {
					t.Fatalf("Expected to get 'updated' for barcode %v but instead got %v\n", v.Barcode, v.Status)
				}
			} else {
				if v.Status != "ignored" {
					t.Fatalf("Expected to get 'ignored' for barcode %v but instead got %v\n", v.Barcode, v.Status)
				}
			}

		}
	}
}

// Check for http 400 when standardprice not provided (null or zero)
func TestPostPricingInvalid(t *testing.T) {

	product.ConnectToDatabase()

	// Switch to test mode to turn off verbose logging
	gin.SetMode(gin.TestMode)

	// Setup our router, just like we did in the main function, and
	// register our routes
	r := gin.Default()
	r.POST("/pricing", postPricing)

	var jsonData = []byte(`[
		{
		  "barcode": "9313820004501",
		  "name": "7Up Lemonade Soda, 1.25 Litre",
		  "standardprice": 0,
		  "specialprice": 2.5
		}
	   ]`)

	// Create the mock request that we would like to test.
	req, err := http.NewRequest(http.MethodPost, "/pricing", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check for HTTP 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, w.Code)
	}
}

// Test GET /products
func TestGetProducts(t *testing.T) {

	product.ConnectToDatabase()

	// Switch to test mode to turn off verbose logging
	gin.SetMode(gin.TestMode)

	// Setup our router, just like we did in the main function, and
	// register our routes
	r := gin.Default()
	r.GET("/products", getProducts)

	// Create the mock request that we would like to test.
	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// New response recorder so we can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check for HTTP 200
	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	} else {
		responseData, _ := io.ReadAll(w.Body)

		var products []product.Product
		json.Unmarshal(responseData, &products)

		//Check for the new updated price of Coca Cola Diet Coke (barcode 9300675009775)
		gotStandardPrice := products[findIndexByBarCode(products, "9300675009775")].StandardPrice
		gotStandardPriceString := strconv.Itoa(int(gotStandardPrice))
		if gotStandardPriceString != updatedPrice {
			t.Fatalf("Expected to get StandardPrice %v for 9300675009775 but instead got %v\n", updatedPrice, gotStandardPriceString)
		}

	}
}

// Search for a product in a products array by barcode and return its index
func findIndexByBarCode(products []product.Product, barcode string) int {
	for idx, product := range products {
		if product.Barcode == barcode {
			return idx
		}
	}
	return -1
}

// Gets a new unique price value. Test helper function
func getNextUniquePrice() int {

	fileLoc := "../nextprice.txt"

	b, errRead := os.ReadFile(fileLoc)
	if errRead != nil {
		panic(errRead)
	}
	val, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	val = val + 1
	if val > 100 {
		val = 1
	}
	errWrite := os.WriteFile(fileLoc, []byte(strconv.Itoa(val)), 0644)
	if errWrite != nil {
		panic(errWrite)
	}
	return val
}
