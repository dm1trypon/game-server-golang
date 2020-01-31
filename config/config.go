package config

import (
	"encoding/json"
	"io/ioutil"
)

// GameConfig - main config of data
var GameConfig Config

func getDataConfigFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// SetConfig method set config data from config's file
func SetConfig(path string) error {
	data, err := getDataConfigFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &GameConfig)
	if err != nil {
		return err
	}

	return nil
}
