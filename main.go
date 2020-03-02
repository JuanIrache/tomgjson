package main

import (
	"fmt"
	"io/ioutil"
	"log"

	utils "github.com/juanirache/tomgjson/utils"
)

func check(e error) {
	if e != nil {
		log.Panic("Error:", e)
	}
}

func main() {
	src, err := ioutil.ReadFile("./samples/only-data.csv")
	check(err)

	converted := utils.ReadCSV(src, 25.0)

	fmt.Println(converted)

}
