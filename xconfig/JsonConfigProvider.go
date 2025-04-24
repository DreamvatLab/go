package xconfig

import (
	"encoding/json"
	"os"

	"github.com/DreamvatLab/go/xjson"
)

type JsonConfigProvider struct {
	RawJson []byte
	MapConfiguration
}

// NewJsonConfigProvider creates a new JsonConfigProvider.
//
// args[0]: The path to the JSON file to read. If not provided, "configs.json" will be used.
//
// Returns a new JsonConfigProvider and an error if the creation fails.
func NewJsonConfigProvider(args ...string) IConfigProvider {
	r := new(JsonConfigProvider)
	r.MapConfiguration = make(MapConfiguration)

	var configFile string
	if len(args) == 0 {
		configFile = "configs.json"
	} else {
		configFile = args[0]
	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	r.RawJson = configData
	err = json.Unmarshal(configData, &r.MapConfiguration)
	if err != nil {
		panic(err)
	}

	return r
}

func (x *JsonConfigProvider) GetStruct(key string, target interface{}) error {
	return xjson.UnmarshalSection(x.RawJson, key, target)
}
