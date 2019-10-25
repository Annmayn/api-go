package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//**************
//todo:
//Combine tables and add a field "Table string" that determines which table to send to
//***************

// //Table1 : First table structure
// type Table1 struct {
// 	Username string `json:username`
// 	Password string `json:password`
// 	// database string
// }

// //Table2 : Second Table structure
// type Table2 struct {
// 	Username string `json:username`
// 	ID       int    `json:id`
// }

//*********useless********
// var m = map[string]*struct{}{"table1": Table1, "table2": &Table2{}}
//*****************

// //Table : interface
// type Table interface{
// 	getDatabase() string
// }

// func (t Table) getDatabase() string{
// 	return t.database
// }

var router *gin.Engine

//check if methods are allowed
//todo: single code to handle all instead of same code for different methods using map
func customHandler(c *gin.Context) {
	var database map[string]interface{}
	jsonFile, _ := os.Open("database.json")
	defer jsonFile.Close()

	byteFile, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteFile), &database)

	switch c.Request.Method {
	case "GET":
		//split by "/"
		url := strings.Split(c.Request.URL.Path, "/")[1:] //ignore the first whitespace as it starts with "/"
		//table name is the 3rd folder (index 3 because url starts with "/") in request : api/v1/[table_name]
		databaseName := url[2]
		//get all the methods allowed in the 'databaseName' database
		db := database[databaseName].(map[string]interface{})
		methodString := db["methods"].(string)
		methods := strings.Split(methodString, ",")

		//assume method is not allowed
		exists := false
		for _, method := range methods {
			if method == "GET" {
				//flag set to indicate GET is allowed
				exists = true
				break
			}
		}
		if exists == false { //if GET call not allowed
			notAllowed(c, "GET")
		} else { //GET call allowed
			handleGet(c, databaseName)
		}

	case "POST":
		url := strings.Split(c.Request.URL.Path, "/")[1:]
		//table name is the 3rd folder (index 2) in request : api/v1/[table_name]
		databaseName := url[2]
		db := database[databaseName].(map[string]interface{})
		methodString := db["methods"].(string)
		methods := strings.Split(methodString, ",")
		exists := false
		for _, method := range methods {
			if method == "POST" {
				exists = true
				break
			}
		}
		if exists == false {
			notAllowed(c, "POST")
		} else {
			//todo: add databaseName (string) as argument
			handleVerification(c)
		}

	default:
		notAllowed(c, "")
	}
}

func handleVerification(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Header("Allow", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "origin, content-type, accept")
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	} else if c.Request.Method == "POST" {
		// c.BindJSON(&u)

		// loc, _ := c.GetQuery("user")

		c.JSON(http.StatusOK, gin.H{
			"user": c.Param("user"),
			"pass": "b",
		})
	}
}

func handleGet(c *gin.Context, databaseName string) {
	//todo: run get dynamically for every get request by iterating one level at a time
	message, _ := c.GetQuery("user")

	//read values from table1.json
	var database map[string]interface{}
	jsonFile, _ := os.Open("database.json")
	defer jsonFile.Close()

	byteFile, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteFile), &database)

	c.String(http.StatusOK, c.Request.URL.Path)
	c.String(http.StatusOK, "Get works!!! "+message+databaseName)
}

func notAllowed(c *gin.Context, method string) {
	if method == "" {
		c.String(http.StatusBadRequest, "Not allowed")
	} else {
		c.String(http.StatusBadRequest, method+" not allowed")
	}
}

func initializeRoutes(r *gin.RouterGroup) {
	//r = localhost:[port number]/api/v1
	r.POST("/*any", customHandler)
	r.OPTIONS("/*any", customHandler)
	r.GET("/*any", customHandler)
	// router.POST("/api/v1", handleVerification)
	// router.OPTIONS("/api/v1", handleVerification)
	// router.GET("/api/v1", handleGet)

	router.Run(":8080")
}

func main() {
	router = gin.Default()
	r1 := router.Group("/api/v1")
	initializeRoutes(r1)
	router.Run()
}
