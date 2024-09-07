package jsutils

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

func ToJsonValue(v interface{}) js.Value {
	json, err := ToJsonMarshal(v)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	return js.ValueOf(json)
}

// converts all exported members of a struct to a map with lowercase field names
func ToJsonMarshal(v interface{}) (JSON, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointers by dereferencing
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("nil pointer received")
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// Check if the kind of value is a supported type
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Slice && val.Kind() != reflect.Map {
		return nil, fmt.Errorf("input must be a struct, slice, or map, received: '%v'", val.Kind())
	}

	data := marshalWithLowercaseNames(val)
	return data, nil
}

// Helper function to recursively marshal a struct to a map with lowercase field names
func marshalWithLowercaseNames(v reflect.Value) JSON {
	result := make(JSON)

	if v.Kind() == reflect.Struct {
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := v.Field(i)

			// Ensure the field is exported and is not a method
			if field.PkgPath == "" && fieldValue.Kind() != reflect.Func {
				fieldName := field.Name
				lowercaseFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]

				if fieldValue.Kind() == reflect.Struct {
					result[lowercaseFieldName] = marshalWithLowercaseNames(fieldValue)
				} else if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
					result[lowercaseFieldName] = processSlice(fieldValue)
				} else if fieldValue.Kind() == reflect.Map {
					result[lowercaseFieldName] = processMap(fieldValue)
				} else {
					result[lowercaseFieldName] = fieldValue.Interface()
				}
			}
		}
	}
	return result
}

// Helper function to process slices
func processSlice(v reflect.Value) interface{} {
	var result []interface{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Struct {
			result = append(result, marshalWithLowercaseNames(elem))
		} else if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
			result = append(result, processSlice(elem))
		} else if elem.Kind() == reflect.Map {
			result = append(result, processMap(elem))
		} else {
			result = append(result, elem.Interface())
		}
	}
	return result
}

// Helper function to process maps
func processMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		if val.Kind() == reflect.Struct {
			result[key.String()] = marshalWithLowercaseNames(val)
		} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			result[key.String()] = processSlice(val)
		} else if val.Kind() == reflect.Map {
			result[key.String()] = processMap(val)
		} else {
			result[key.String()] = val.Interface()
		}
	}
	return result
}
