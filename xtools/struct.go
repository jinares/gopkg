package xtools

import (
	"html"
	"reflect"
	"strconv"
)

// MapToStruct MapToStruct
func MapToStruct(data map[string]string, out interface{}) error {

	ss := reflect.ValueOf(out).Elem()
	for i := 0; i < ss.NumField(); i++ {
		val, isok := data[ss.Type().Field(i).Tag.Get("json")]
		if isok == false {
			continue
		}
		name := ss.Type().Field(i).Name

		switch ss.Field(i).Kind() {
		case reflect.String:
			ss.FieldByName(name).SetString(val)
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetInt(int64(i))
		case reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.Atoi(val)
			if err != nil {

				continue
			}
			ss.FieldByName(name).SetUint(uint64(i))
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetFloat(f)
		default:

		}
	}
	return nil
}

//StructToMap StructToMap
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//StructToMapByTag StructToMapByTag
func StructToMapByTag(obj interface{}, tag string) map[string]interface{} {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if tag == "" {
		tag = "json"
	}
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get(tag)
		if tagName == "" {
			continue
		}
		data[tagName] = v.Field(i).Interface()
	}
	return data
}

//XSSHandleStruct 对struct的指定字段进行xss过滤
func XSSHandleStruct(data interface{}, tag string, field []string) error {
	ss := reflect.ValueOf(data).Elem()

	opMap := map[string]interface{}{}
	for _, v := range field {
		opMap[v] = nil
	}

	for i := 0; i < ss.NumField(); i++ {
		tag := ss.Type().Field(i).Tag.Get(tag)
		if len(tag) <= 0 {
			continue
		}
		if _, ok := opMap[tag]; !ok {
			continue
		}

		name := ss.Type().Field(i).Name

		if val, ok := ss.FieldByName(name).Interface().(string); ok {
			xSSHandleString(&val)
			ss.FieldByName(name).SetString(val)
		}

	}

	return nil
}

//XSSHandleString 对字符串进行xss过滤
func xSSHandleString(data *string) error {

	*data = html.EscapeString(*data)
	return nil
}
