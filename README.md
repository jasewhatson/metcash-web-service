## Usage instructions / How to run

The project has been developed in *go version 1.19.1* however the code base should be compatible with older versions of Golang

All below commands are to be run from within the src of the project folder (eg *cd /Users/JasonWhatson/src/metcash-web-service/src*)

**How to install the package dependencies**

The following 3rd party packages are required 

[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) - Provides REST API

[github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - Provides data persistence via Sqlite

These can be installed via running the following from the src directory

> go get

**How to run the unit tests located in** *main_test.go*

The following tests scenarios _TestPostPricing_, _TestGetProducts_ and _TestPostPricingInvalid_ are provided in *main_test.go*

> go test  
>
or
>
> go test -v

**How to run the main web service**

We can start our restful web service on port 8880 via the following command

> go run main.go

**How to make test requests to our web service**

Alternative to running the unit tests above. We can manually test the web service via the following

> curl --request "GET" http://localhost:8880/products

>curl http://localhost:8880/pricing --include --header "Content-Type: application/json" --request "POST" --data-binary "@../Prices.json"

## Documentation

The requirements do not specify what to do on _POST /pricing_ if an unknown product (one which is not in the DB) is in the request. So, for this scenario, I return _status = "notfound"_ for this product. Alternatively, we could just silently ignore this product from the update.

Data persistence is provided via [SQLite](https://en.wikipedia.org/wiki/SQLite). With the database file for our project called _products.db_

**Error and failure handling** 

In our product package errors are handled by checking the result of appropriate calls to things like DB reads/writes via calling checkErr(). If an error is found, panic and print the error msg & stack trace. Alternatively, we good log the error here. 

**Validation** 

Validation for the _/pricing_ JSON payload is done via the _gin_ package and will result in an 'HTTP/1.1 400 Bad Request' response if there are any errors in the provided JSON payload such as syntax or wrong data type. Also if the value for _standardprice_ is provided as null or zero, 'HTTP/1.1 400 Bad Request' is returned. This is tested in the unit test _TestPostPricingInvalid_ 

## Need clarification on requirements 

From the requirements doc, it says _'If a price already exists for this product and is not in the API call it should be deleted'_

I am not sure what is meant by this. It should be deleted from where? 
If this is required in my implementation, can you please clarify what is meant by this so I can add it to the project? Thanks
