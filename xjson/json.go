package xjson

import (
	"encoding/json"

	"github.com/Lukiya/go/xbytes"
	"github.com/Lukiya/go/xerr"
	"github.com/tidwall/gjson"
)

// UnmarshalSection unmarshals the section of the JSON data.
//
// data: The JSON data to unmarshal.
// path: The path to the section to unmarshal.
// target: The target to unmarshal the section to.
//
// Returns an error if the unmarshalling fails.
func UnmarshalSection(data []byte, path string, target interface{}) error {
	v := gjson.GetBytes(data, path)
	err := json.Unmarshal(xbytes.StrToBytes(v.Raw), target)
	if err != nil {
		return xerr.WithStack(err)
	}
	return nil
}
