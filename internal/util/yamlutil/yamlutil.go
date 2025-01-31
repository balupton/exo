package yamlutil

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// Prefers 2 space indentation.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

func MarshalString(v interface{}) (string, error) {
	bs, err := Marshal(v)
	return string(bs), err
}

func MustMarshalString(v interface{}) string {
	s, err := MarshalString(v)
	if err != nil {
		panic(err)
	}
	return s
}

func UnmarshalString(s string, v interface{}) error {
	return yaml.Unmarshal([]byte(s), v)
}

func MustUnmarshalString(s string, v interface{}) {
	if err := UnmarshalString(s, v); err != nil {
		panic(err)
	}
}
