package main

import (
	"io/ioutil"
	"encoding/json"
)

/*
	{
		"settings": {
			
		},

		"events": [
			{ 
				"event": "power-menu",
				"command": "rofi ..."
			}
		]
	}
*/

type Configuration struct {
	Settings ConfigSettings	`json:"settings"`
	Events []ConfigEvent `json:"events"`
}

type ConfigSettings struct {
	Lemonbar string 	`json:"lemonbar"`
	Font string 		`json:"font"`
}

type ConfigEvent struct {
	Event string 	`json:"event"`
	Command string	`json:"command"`
}

func LoadConfig(path string) Configuration {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	config := Configuration{}
	json.Unmarshal(content, &config)

	return config
}
