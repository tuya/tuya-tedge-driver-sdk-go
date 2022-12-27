package utils

import (
	"bytes"
	"encoding/json"
)

func JsonDecoder(payload []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(payload))
	d.UseNumber()
	return d.Decode(v)
}
