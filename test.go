package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Tablex struct {
	Url interface{} `json:table1`
	Uri []string    `json:table2`
}

func main() {
	jsonFile, _ := os.Open("product.json")
	defer jsonFile.Close()
	// fmt.Println(jsonFile)
	// var t Table1
	var t map[string]interface{}
	byteFile, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteFile), &t)
	// fmt.Println(string(byteFile))
	fmt.Println(t.table1)
}
