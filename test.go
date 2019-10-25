package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Book struct {
	Id      int
	Title   string
	Price   float32
	Authors []string
}

// type Tablex struct {
// 	Url interface{} `json:table1`
// 	Uri []string    `json:table2`
// }

//------> c.Request.URL.Path -> request url

func main() {
	jsonFile, _ := os.Open("product.json")
	defer jsonFile.Close()
	// fmt.Println(jsonFile)
	// var t Table1
	var t map[string]interface{}
	byteFile, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteFile), &t)
	// fmt.Println(string(byteFile))
	// t = t.(map[string]interface{})
	// fmt.Printf("%T", t)
	for k, _ := range t {
		tmp, _ := t[k].(map[string]interface{})
		fmt.Println(tmp)
	}
	// fmt.Println(t.table1)

	//*****************Reading struct elements*****************
	// a := "localhost:8080/api/v1/key"
	// b := strings.Split(a, "/")
	// fmt.Println(b)
	// book := Book{}
	// e := reflect.ValueOf(&book).Elem()

	// //get struct elements
	// for i := 0; i < e.NumField(); i++ {
	// 	varName := e.Type().Field(i).Name
	// 	varType := e.Type().Field(i).Type
	// 	varValue := e.Field(i).Interface()
	// 	fmt.Printf("%v %v %v\n", varName, varType, varValue)
	// }
}
