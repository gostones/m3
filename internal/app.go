package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AppDesc struct {
	Name             string `json:"name"`
	Command          string `json:"command"`
	AutoRestart      bool   `json:"autoRestart"`
	After            string `json:"after"`
	WorkingDirectory string `json:"workDir"`
}

//LoadAppDesc loads app configuration by name from ./, $dhnt_base/etc/ or $HOME/dhnt/etc/
func LoadAppDesc(name string) (*AppDesc, error) {
	// current dir
	desc, err := loadJSON(name)
	if err == nil {
		return desc, nil
	}

	// base := GetDefaultBase()
	// if base != "" {
	// 	desc, err = loadJSON(fmt.Sprintf("%v/etc/%v.json", base, name))
	// 	if err == nil {
	// 		return desc, nil
	// 	}
	// }

	return nil, fmt.Errorf("configuration not found for %v", name)
}

func loadJSON(file string) (*AppDesc, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var desc AppDesc
	err = json.Unmarshal(data, &desc)
	return &desc, err
}
