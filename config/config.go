package config 

import (
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
	Settings ConfigSettings	`json:"settings"`
	Template ConfigTemplate `json:"template"`
	Blocks map[string] ConfigBlock `json:"blocks"`
	Buttons map[string] ConfigButton `json:"buttons"`
}

type ConfigTemplate struct {
	Left []string `json:"left"`
	Center []string `json:"center"`
	Right []string `json:"right"`
}

type ConfigSettings struct {
	Lemonbar string 	`json:"lemonbar"`
	Font string 		`json:"font"`
}

type ConfigBlock struct {
	Name string `json:"name"`
	Module string `json:"module"`
	Command string `json:"command"`
	Template string `json:"template"`
	Interval string `json:"interval"`
}

type ConfigButton struct {
	Name string `json:"name"`
	Command string `json:"command"`
	Label string `json:"label"`
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
