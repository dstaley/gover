package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func configure() {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		println(`ERROR: conf.json file not found. Please run "gover setup" to create one.`)
		os.Exit(-1)
	} else {
		blob, err := ioutil.ReadFile("./config.json")
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(blob, &config)
	}
}
