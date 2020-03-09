package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dm1trypon/game-server/servicedata"
	"github.com/ivahaev/go-logger"
	"github.com/qri-io/jsonschema"
)

// LC - Logging category
const LC = "[Config] >> "

// IsValidConfig - a method that checks the config for validation.
// If the config is correct, returns true, otherwise - false.
func IsValidConfig(pathConfig string, pathSchema string) bool {
	configData, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		logger.Error(LC + err.Error())
		return false
	}

	schemaData, err := ioutil.ReadFile(pathSchema)
	if err != nil {
		logger.Error(LC + err.Error())
		return false
	}

	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaData, rs); err != nil {
		logger.Error(LC + err.Error())
		return false
	}

	if errors, err := rs.ValidateBytes(configData); len(errors) > 0 {
		for _, valErr := range errors {
			logger.Error(LC + valErr.Error())
		}

		return false
	} else if err != nil {
		logger.Error(LC + err.Error())
		return false
	}

	if err := json.Unmarshal(configData, &servicedata.GameConfig); err != nil {
		logger.Error(LC + err.Error())
		return false
	}

	return true
}
