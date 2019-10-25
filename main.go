package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Database : custom database
var Database = map[string]map[string]string{
	"table1": {
		"url":     "localhost:8080/api/v1/",
		"uri":     "table1/",
		"methods": "get,post",
	},
	"table2": {
		"url":     "localhost:8080/api/v1/",
		"uri":     "table2/",
		"methods": "get",
	},
}

//Table1 : First table
type Table1 struct {
	Username string `json:username`
	Password string `json:password`
}

//Table2 : Second Table
type Table2 struct {
	Username string `json:username`
	ID       int    `json:id`
}

var router *gin.Engine

func initializeRoutes() {
	router.POST("/api/v1", handleVerification)
	router.OPTIONS("/api/v1", handleVerification)
	router.GET("/api/v1", handleGet)

	router.Run(":8080")
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
	c.String(http.StatusOK, "Get works!!! "+message)
}

func main() {
	router = gin.Default()
	initializeRoutes()
	router.Run()
}
