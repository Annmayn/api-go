package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

//**************
//todo:
//Combine tables and add a field "Table string" that determines which table to send to
//***************

//Table1 : First table structure
type Table1 struct {
	Username string `json:"username"`
	Password string `json:"passcode"`
	// database string
}

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

		table := reflect.ValueOf(database[databaseName])
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
			col := reflect.ValueOf("schema")
			handlePost(c, url[2:], table.MapIndex(col))
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

func handlePost(c *gin.Context, dir []string, schema reflect.Value) {
	fmt.Println("Hello")
	// var t1 Table1
	// c.BindJSON(&t1)
	// fmt.Printf("T1: %+v", t1)
	//*********************************

	//*****************************read post body using post-form*******************************
	// loc, _ := c.GetQuery("user")
	// url := strings.Split(c.Request.URL.Path, "/")[3:] //skip the first whitespace due to trailing '/', "api" and "v1"

	// c.Request.ParseForm()
	// for k, v := range c.Request.PostForm {
	// 	fmt.Println(k, ": ", v)
	// }

	//todo: kv could be map[reflect.Value]interface{}
	var kv map[string]interface{}
	//dynamically create a map to check for valid post request
	for _, k := range schema.MapKeys() {
		kv[k.Interface().(string)] = ""
	}
	fmt.Println(kv)

	//*************read post body using interface*********************
	var jsonInterface interface{}
	c.BindJSON(&jsonInterface)
	fmt.Println(jsonInterface)
	reflectJSON := reflect.ValueOf(jsonInterface)
	switch reflectJSON.Kind() {
	case reflect.Map:
		for _, k := range reflectJSON.MapKeys() {
			// fmt.Println(k, reflectJSON.MapIndex(k))
			if _, ok := kv[k.Interface().(string)]; ok {
				kv[k.Interface().(string)] = reflectJSON.MapIndex(k)
			}
			// z := reflect.ValueOf("passcode")
			// fmt.Println(k, v.MapIndex(z))
		}
	}
	// fmt.Printf("%T\n", jsonInterface)
	//*********************************

	// fo, _ := os.Open("table2.json")
	// defer fo.Close()

	// w := bufio.NewWriter(fo)
	// _, _ = w.Write(jsonInterface.([]byte))

	fmt.Println(json.Marshal(kv))
	j, _ := json.Marshal(kv)
	// jsonByte := []byte(fmt.Sprintf("%v", jsonInterface.(map[string]interface{})))
	_ = ioutil.WriteFile(dir[0]+".json", j, 0777)

	//**********************************

	file, _ := os.Open(dir[0] + ".json")
	defer file.Close()
	content, _ := ioutil.ReadAll(file)

	// b := c.PostForm("val")
	fmt.Println(string(content))
	// fmt.Println(b["a"])

	// c.JSON(http.StatusOK, gin.H{
	// 	"user": c.Param("user"),
	// 	"pass": "b",
	// })
	c.String(http.StatusOK, "")
}

func handleGet(c *gin.Context, dir []string) {
	//todo: run get dynamically for every get request by iterating one level at a time
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
