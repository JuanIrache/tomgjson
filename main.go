package main

import (
	"encoding/json"
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
	mgjson := utils.FormatMgjson(converted, "github.com/juanirache/tomgjson")

	f, err := os.Create("./out.json")
	check(err)

	defer f.Close()

	doc, err := json.Marshal(mgjson)
	check(err)

	_, err = f.Write(doc)
	check(err)

}
