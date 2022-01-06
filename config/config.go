package config 

import (
	"io/ioutil"
	"encoding/json"
)

type Configuration struct {
	Settings ConfigSettings	`json:"settings"`
	Template ConfigTemplate `json:"template"`
	Blocks map[string] ConfigBlock `json:"blocks"`
}

type ConfigTemplate struct {
	Left []string `json:"left"`
	Center []string `json:"center"`
	Right []string `json:"right"`
}

type ConfigSettings struct {
	Lemonbar string 	`json:"lemonbar"`
	Font string 		`json:"font"`
	Separator string 	`json:"separator"`
}

type ConfigBlock struct {
	Module string `json:"module"`
	Template string `json:"template"`
	Interval string `json:"interval"`
	OnClick string `json:"on-click"`
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
