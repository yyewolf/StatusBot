//go:generate goversioninfo
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func saveData() error {
	//File isn't openable
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		return errors.New("cannot save/create db file")
	}
	return nil
}

func loadData() {
	//Read file
	file, err := ioutil.ReadFile("data.json")

	//Create if it doesn't exists
	if err != nil {
		err := saveData()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	//Try to put the file in IPS
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		err := saveData()
		if err != nil {
		}
	}
}
