package gocache

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	data := "hello"
	encoded, err := encode(data)
	if err != nil {
		t.Error(err)
		return
	}
	expected := []byte{34, 104, 101, 108, 108, 111, 34, 10}
	if !reflect.DeepEqual(encoded, expected) {
		t.Error("unexpected encoded result")
		return
	}
}

func TestDecode(t *testing.T) {
	data := []byte{34, 104, 101, 108, 108, 111, 34, 10}
	str := ""
	err := decode(data, &str)
	if err != nil {
		t.Error(err)
		return
	}
	expected := "hello"
	if str != expected {
		t.Error("unexpected decoded result")
		return
	}
}
