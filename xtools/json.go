package xtools

import (
	"bytes"
	"encoding/json"
	"strings"
)

//JSONToStr JsonToStr
func JSONToStr(dstr interface{}) string {
	switch vv := dstr.(type) {
	case string:
		return vv
	}
	data, err := json.Marshal(dstr)
	if err != nil {
		return ""
	}
	return string(data)

}

//JSONMarshal JsonToStr
func JSONMarshal(dstr interface{}) string {

	switch vv := dstr.(type) {
	case string:
		return vv
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(dstr)

	return strings.TrimSpace(string(buffer.Bytes()))

}

//ToJson ToJson
func ToJSON(dstr string, data interface{}) error {
	err := json.Unmarshal([]byte(dstr), &data)
	return err
}
