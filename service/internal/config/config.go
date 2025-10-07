package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
)

func NewNews() NewsServer {
	server := NewsServer{}
	buildconfig(&server, configPath)
	return server
}

const (
	configPath = "./config.json"
)

func buildconfig(cfg interface{}, configPath string) {

	file, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open config file: %w", err))
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %w", err))
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		panic(fmt.Sprintf("failed to unmarshal json: %w", err))
	}

	decoderConfigOption := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
		),
		Result: cfg,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfigOption)
	if err != nil {
		panic(fmt.Sprintf("Error creating decoder: %v", err))
	}

	if err = decoder.Decode(raw); err != nil {
		panic(fmt.Sprintf("Error decoding config: %v", err))
	}

}
