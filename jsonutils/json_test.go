package jsonutils

import (
	"strings"
	"testing"
)

func TestFromReader_Then_ReturnObject(t *testing.T) {

	type CarModel struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Car  string `json:"car"`
	}

	reader := strings.NewReader(`{"name":"John","age":30,"car":null}`)

	carModel, err := FromReader[CarModel](reader)
	if err != nil {
		t.Error(err)
	}

	if carModel.Name != "John" {
		t.Error("Name is not John")
	}

	if carModel.Age != 30 {
		t.Error("Age is not 30")
	}

	if carModel.Car != "" {
		t.Error("Car is not empty")
	}
}

func TestFromReader_Then_ReturnString(t *testing.T) {

	reader := strings.NewReader(`{"name":"John","age":30,"car":null}`)

	str, err := FromReader[string](reader)
	if err != nil {
		t.Error(err)
	}

	if str != `{"name":"John","age":30,"car":null}` {
		t.Error("String is not correct")
	}
}

func TestFromReader_Then_ReturnBytes(t *testing.T) {

	reader := strings.NewReader(`{"name":"John","age":30,"car":null}`)

	bytes, err := FromReader[[]byte](reader)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != `{"name":"John","age":30,"car":null}` {
		t.Error("Bytes is not correct")
	}
}

func TestFromString_Then_ReturnObject(t *testing.T) {

	type CarModel struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Car  string `json:"car"`
	}
	
	carModel, err := FromString[CarModel](`{"name":"John","age":30,"car":null}`)
	if err != nil {
		t.Error(err)
	}

	if carModel.Name != "John" {
		t.Error("Name is not John")
	}

	if carModel.Age != 30 {
		t.Error("Age is not 30")
	}

	if carModel.Car != "" {
		t.Error("Car is not empty")
	}
}
