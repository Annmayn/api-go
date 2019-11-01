package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// )

// // type Book struct {
// // 	Id      int
// // 	Title   string
// // 	Price   float32
// // 	Authors []string
// // }

// // type Tablex struct {
// // 	Url interface{} `json:table1`
// // 	Uri []string    `json:table2`
// // }

// //------> c.Request.URL.Path -> request url

// func main() {
// 	// jsonFile, _ := os.Open("database.json")
// 	// defer jsonFile.Close()

// 	//******typecasting nested json***************

// 	// var t map[string]interface{}
// 	// byteFile, _ := ioutil.ReadAll(jsonFile)
// 	// json.Unmarshal(byteFile, &t)
// 	// fmt.Println(t["table1"])
// 	// a := t["table1"].(map[string]interface{})
// 	// c := strings.Split(a["methods"].(string), ",")
// 	// fmt.Println(c)
// 	// for k := range t {
// 	// 	tmp, _ := t[k].(map[string]interface{})
// 	// 	fmt.Println(tmp)
// 	// }
// 	//*********************************************

// 	//*****************Reading struct elements*****************
// 	// a := "localhost:8080/api/v1/key"
// 	// b := strings.Split(a, "/")
// 	// fmt.Println(b)
// 	// book := Book{}
// 	// e := reflect.ValueOf(&book).Elem()

// 	// //get struct elements
// 	// for i := 0; i < e.NumField(); i++ {
// 	// 	varName := e.Type().Field(i).Name
// 	// 	varType := e.Type().Field(i).Type
// 	// 	varValue := e.Field(i).Interface()
// 	// 	fmt.Printf("%v %v %v\n", varName, varType, varValue)
// 	// }

// 	// a := "/api/key/v1"
// 	// a = a[1:]
// 	// fmt.Println(a)

// 	var database map[string]interface{}
// 	jsonFile, _ := os.Open("database.json")
// 	defer jsonFile.Close()
// 	byteFile, _ := ioutil.ReadAll(jsonFile)
// 	json.Unmarshal([]byte(byteFile), &database)

// 	// for k, v := range database {
// 	// 	fmt.Println(k, ":")
// 	// 	z := reflect.ValueOf(v)
// 	// 	for _, x := range z.MapKeys() {
// 	// 		fmt.Printf("%T\n", z.MapIndex(x))
// 	// 		fmt.Println("-> ", x, " : ", z.MapIndex(x))
// 	// 	}
// 	// 	// fmt.Println(z.MapKeys())
// 	// }

// 	a := map[string]string{"apple": "1"}
// 	if _, ok := a["ball"]; !ok {
// 		fmt.Println("BAll doesn't exist")
// 	}
// }
