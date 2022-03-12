package gocache

import (
	"bytes"
	"encoding/json"
)

// encode encodes any Go type to a slice of bytes
func encode(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(data)
	return buf.Bytes(), err
}

// decode decodes a slice of byte (generated with the encode func) into a Go type that it stores in the provided pointer
func decode(data []byte, intoptr interface{}) error {
	buf := bytes.NewBuffer(data)
	return json.NewDecoder(buf).Decode(intoptr)
}
