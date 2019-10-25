package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//Database : custom database
var Database = map[string]map[string]string{
	"table1": {
		"url":     "localhost:8080/api/v1/table1",
		"uri":     "/",
		"methods": "POST,GET",
	},
	"table2": {
		"url":     "localhost:8080/api/v1/table2",
		"uri":     "/",
		"methods": "GET",
	},
}

//**************
//todo:
//Combine tables and add a field "Table string" that determines which table to send to
//***************

//Table1 : First table structure
type Table1 struct {
	Username string `json:username`
	Password string `json:password`
	// database string
}

//Table2 : Second Table structure
type Table2 struct {
	Username string `json:username`
	ID       int    `json:id`
}

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
	switch c.Request.Method {
	case "GET":
		//split by "/"
		url := strings.Split(c.Request.URL.Path, "/")
		//table name is the 3rd folder (index 3 because url starts with "/") in request : api/v1/[table_name]
		databaseName := url[3]
		//get all the methods allowed in the 'databaseName' database
		methods := strings.Split(Database[databaseName]["methods"], ",")
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
			handleGet(c)
		}

	case "POST":
		url := strings.Split(c.Request.URL.Path, "/")
		//table name is the 3rd folder (index 2) in request : api/v1/[table_name]
		tableName := url[3]
		methods := strings.Split(Database[tableName]["methods"], ",")
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
		var u Table1
		c.BindJSON(&u)
		c.JSON(http.StatusOK, gin.H{
			"user": u.Username,
			"pass": u.Password,
		})
	}
}

func handleGet(c *gin.Context) {
	message, _ := c.GetQuery("m")
	c.String(http.StatusOK, c.Request.URL.Path)
	c.String(http.StatusOK, "Get works!!! "+message)
}

func notAllowed(c *gin.Context, method string) {
	if method == "" {
		c.String(http.StatusBadRequest, "Not allowed")
	} else {
		c.String(http.StatusBadRequest, method+" not allowed")
	}
}

func initializeRoutes(r *gin.RouterGroup) {
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
