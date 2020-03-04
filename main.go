package main

import (
	"io/ioutil"
	"log"
	"os"

	utils "github.com/juanirache/tomgjson/utils"
)

func check(e error) {
	if e != nil {
		log.Panic("Error:", e)
	}
}

func main() {
	src, err := ioutil.ReadFile("./samples/multiple-data.csv")
	check(err)

	converted := utils.ReadCSV(src, 25.0)
	doc := utils.FormatMgjson(converted, "github.com/juanirache/tomgjson")

	f, err := os.Create("./out.json")
	check(err)

	defer f.Close()

	_, err = f.Write(doc)
	check(err)

}
