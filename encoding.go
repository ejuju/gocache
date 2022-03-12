package gdb

import (
	"bytes"
	"encoding/json"
)

func encode(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(data)
	return buf.Bytes(), err
}

// decode data into pointer to struct, map, list or whatever...
func decode(data []byte, intoptr interface{}) error {
	buf := bytes.NewBuffer(data)
	return json.NewDecoder(buf).Decode(intoptr)
}
