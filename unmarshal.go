package gokhipu

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func unmarshalJSON(r io.Reader, v interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
