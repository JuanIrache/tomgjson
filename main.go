package main

import (
	"fmt"
	"io/ioutil"
	"log"

	readcsv "github.com/juanirache/tomgjson/utils"
)

func check(e error) {
	if e != nil {
		log.Panic("Error:", e)
	}
}

func main() {
	src, err := ioutil.ReadFile("./samples/labelled-data.csv")
	check(err)

	converted := readcsv.ReadCSV(src)

	fmt.Println(converted)

}
