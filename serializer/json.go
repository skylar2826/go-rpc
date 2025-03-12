package serializer

import (
	json2 "encoding/json"
)

type Json struct {
}

func (j *Json) Code() uint8 {
	return 1
}

func (j *Json) Encode(data interface{}) ([]byte, error) {
	return json2.Marshal(data)
}

func (j *Json) DeCode(bs []byte, data interface{}) error {
	return json2.Unmarshal(bs, data)
}
