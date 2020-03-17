package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {

	byt := []byte(`{"num":"6.13","strs":["a","b"]}`)
	var dat map[string]interface{}

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	if dat["num"] == nil {
		fmt.Println("Is nil")
		return
	}

	if reflect.TypeOf(dat["num"]).String() == "string" {
		fmt.Println("Is string")
		return
	}

	fmt.Println(dat["num"].(string))
}
