package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

//todo: single code to handle all request instead of same code for different requests (using map)

//check if methods are allowed
func customHandler(c *gin.Context) {
	//read config file from database.json
	var database map[string]interface{}
	jsonFile, _ := os.Open("database.json")
	defer jsonFile.Close()

	byteFile, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteFile), &database)

	switch c.Request.Method {
	case "GET":
		//split by "/"
		url := strings.Split(c.Request.URL.Path, "/")[1:] //ignore the first whitespace as it starts with "/"

		fmt.Println(c.Request.URL.Path)
		fmt.Println(c.Request.URL.Query())

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
			handleGet(c, url[2:])
		}

	case "POST":
		url := strings.Split(c.Request.URL.Path, "/")[1:]
		//table name is the 3rd folder (index 2) in request : api/v1/[table_name]
		databaseName := url[2]

		db := database[databaseName].(map[string]interface{})

		schema := db["schema"].(map[string]interface{})
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
			handlePost(c, url[2:], schema)
		}

	default:
		notAllowed(c, "")
	}
}

// func handleVerification(c *gin.Context) {
// 	if c.Request.Method == "OPTIONS" {
// 		c.Header("Allow", "POST, GET, OPTIONS")
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Headers", "origin, content-type, accept")
// 		c.Header("Content-Type", "application/json")
// 		c.Status(http.StatusOK)
// 	} else if c.Request.Method == "POST" {
// 		// c.BindJSON(&u)

//#region
// 		//*****************************current work*******************************
// 		// loc, _ := c.GetQuery("user")
// 		// url := strings.Split(c.Request.URL.Path, "/")[3:] //skip the first whitespace due to trailing '/', "api" and "v1"
// 		// c.Request.ParseForm()
// 		// for k, v := range c.Request.PostForm {
// 		// 	fmt.Println(k, ": ", v)
// 		// }
// 		b := c.PostFormArray("val")
// 		fmt.Println(b)
// 		// fmt.Println(b["a"])
// 		fmt.Println("Done")
// 		// c.JSON(http.StatusOK, gin.H{
// 		// 	"user": c.Param("user"),
// 		// 	"pass": "b",
// 		// })
// 	}
// }

func validateKV(schema map[string]interface{}, kv map[string]interface{}) []string {
	var missingData []string
	//schema -> map[string]map[string]string
	for k := range schema {
		if _, ok := kv[k]; !ok { //if doesn't exist
			//schemaDetail -> map[string]string
			schemaDetail := schema[k].(map[string]interface{})
			//check if field was required to be present
			if schemaDetail["required"].(string) == "true" {
				missingData = append(missingData, k)
			}
		}
	}
	return missingData
}

func handlePost(c *gin.Context, dir []string, schema map[string]interface{}) {
	kv := make(map[string]interface{})

	var jsonInterface map[string]interface{}
	c.BindJSON(&jsonInterface)

	for k := range jsonInterface {
		//check if key exists
		if _, ok := schema[k]; ok { //if exists in schema, update
			kv[k] = jsonInterface[k]
		}
	}

	//make sure request body is valid
	missingData := validateKV(schema, kv)

	if len(missingData) == 0 { //no missing data, all good
		marshalledJSON, _ := json.Marshal(kv)

		//todo: update file writer to append (currently overwrites)
		_ = ioutil.WriteFile(dir[0]+".json", marshalledJSON, 0777)

		c.String(http.StatusOK, "")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"MissingFields": missingData})
	}
}

func handleGet(c *gin.Context, dir []string) {
	//tip: this can be used in the future
	// message, _ := c.GetQuery("user")

	// //read values from table1.json
	// var database map[string]interface{}
	// jsonFile, _ := os.Open("database.json")
	// defer jsonFile.Close()

	// byteFile, _ := ioutil.ReadAll(jsonFile)
	// json.Unmarshal([]byte(byteFile), &database)

	// c.String(http.StatusOK, c.Request.URL.Path)
	// c.String(http.StatusOK, "Get works!!! "+message+databaseName)

	file, _ := os.Open(dir[0] + ".json")
	content, _ := ioutil.ReadAll(file)
	// fmt.Println(content)
	c.String(http.StatusOK, string(content))
}

func notAllowed(c *gin.Context, method string) {
	if method == "" {
		c.JSON(http.StatusBadRequest, "Not allowed")
	} else {
		c.String(http.StatusBadRequest, method+" not allowed")
	}
}

func initializeRoutes(r *gin.RouterGroup) {
	//r = localhost:[port number]/api/v1
	r.POST("/*any", customHandler)
	r.OPTIONS("/*any", customHandler)
	r.GET("/*any", customHandler)

	router.Run(":8080")
}

func main() {
	router = gin.Default()
	r1 := router.Group("/api/v1")
	initializeRoutes(r1)
	router.Run()
}
