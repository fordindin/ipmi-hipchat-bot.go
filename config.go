package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func handleError(err error, filepath string, fatal bool) {
	if err != nil {
		log.Println("Error reading configuration file", filepath)
		if fatal {
			log.Fatal("error:", err)
		} else {
			log.Println("error:", err)
			log.Println("Continuing with built-in defaults")
		}
	}
}

func readConfig(filepath string) error {
	source, err := ioutil.ReadFile(filepath)
	handleError(err, filepath, false)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(source, &config)
	handleError(err, filepath, true)
	return err
}
