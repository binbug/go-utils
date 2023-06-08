package jsonutils

import (
	"bytes"
	"encoding/json"
	"io"
)

func FromReader[T any](reader io.Reader) (t T, err error) {
	switch any(t).(type) {
	case string:
		data, err := io.ReadAll(reader)
		if err != nil {
			return t, err
		}

		t = any(string(data)).(T)
		return t, nil
	case []byte:
		data, err := io.ReadAll(reader)
		if err != nil {
			return t, err
		}
		t = any(data).(T)
		return t, nil
	default:
		err = json.NewDecoder(reader).Decode(&t)
	}

	return t, err

}

func FromBytes[T any](data []byte) (t T, err error) {

	switch any(t).(type) {
	case string:
		t = any(string(data)).(T)
	case []byte:
		t = any(data).(T)
	default:
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&t)
	}

	return t, err
}

func FromString[T any](data string) (t T, err error) {
	return FromBytes[T]([]byte(data))
}

func ToJSON(o interface{}) string {

	switch o.(type) {
	case string:
		return o.(string)
	default:
		bf := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(bf)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.SetIndent("", "")
		jsonEncoder.Encode(o)
		return bf.String()
	}
}
