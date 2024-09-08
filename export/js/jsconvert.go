package jsconvert

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

type JSON = map[string]any

func ToJsArray(value any) js.Value {
	// Get the reflection value of the input
	refValue := reflect.ValueOf(value)

	// Check if the value is a slice or an array
	if refValue.Kind() != reflect.Slice && refValue.Kind() != reflect.Array {
		panic("provided value is neither a slice nor an array")
	}

	// Create a JavaScript array with the length of the Go slice or array
	jsArray := js.Global().Get("Array").New(refValue.Len())

	// Fill the JavaScript array with the values from the Go slice or array
	for i := 0; i < refValue.Len(); i++ {
		elem := refValue.Index(i).Interface()
		jsArray.SetIndex(i, js.ValueOf(elem))
	}

	return jsArray
}

func ToJsValue(v any) js.Value {
	jsonMap, err := ToJsonMarshal(v)
	if err != nil {
		return js.Global().Get("Error").New(err.Error())
	}
	typ := reflect.TypeOf(reflect.ValueOf(v))
	if typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Struct {
		AddStructMethods(reflect.ValueOf(v), jsonMap)
	}
	return js.ValueOf(jsonMap)
}

func ToJsonMarshal(v any) (JSON, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointers by dereferencing
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("nil pointer received")
		}
		val = val.Elem()
	}

	// Check if the kind of value is a supported type
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Slice && val.Kind() != reflect.Map {
		return nil, fmt.Errorf("input must be a struct, slice, or map (or pointer to one of those), received: '%v'", val.Kind())
	}

	data := marshalToCamelCaseWithMethods(val, true)
	return data, nil
}

func marshalToCamelCaseWithMethods(v reflect.Value, includeMethods bool) JSON {
	result := make(JSON)
	if v.Kind() == reflect.Struct {
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" {
				// Ensure the field is exported
				continue
			}
			fieldValue := v.Field(i)
			fieldName := field.Name
			lowercaseFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]

			if fieldValue.Kind() == reflect.Struct {
				// Recursively process nested structs and their methods
				nestedResult := marshalToCamelCaseWithMethods(fieldValue, true)
				result[lowercaseFieldName] = nestedResult
			} else if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Struct {
				// Handle pointer to struct
				nestedResult := marshalToCamelCaseWithMethods(fieldValue.Elem(), true)
				result[lowercaseFieldName] = nestedResult
			} else if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
				result[lowercaseFieldName] = processSlice(fieldValue)
			} else if fieldValue.Kind() == reflect.Map {
				result[lowercaseFieldName] = processMap(fieldValue)
			} else {
				result[lowercaseFieldName] = fieldValue.Interface()
			}
		}
	}

	// Optionally include methods tied to the struct
	if includeMethods && (v.Kind() == reflect.Ptr || v.Kind() == reflect.Struct) {
		AddStructMethods(v, result)
	}

	return result
}

// Helper function to add registered struct methods to the resulting JSON
func AddStructMethods(v reflect.Value, result JSON) {
	typ := v.Type()
	// Iterate over the methods of the pointer or struct type
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if method.PkgPath != "" {
			// Method must be exported to be added
			continue
		}
		methodName := method.Name
		lowercaseMethodName := strings.ToLower(methodName[:1]) + methodName[1:]
		// Directly pass the method to js.FuncOf
		jsMethod := js.FuncOf(func(this js.Value, args []js.Value) any {
			goMethod := v.MethodByName(methodName)
			// Validate argument count
			if !goMethod.IsValid() || goMethod.Type().NumIn() != len(args) {
				return fmt.Errorf("invalid method or argument count mismatch for %s", methodName)
			}
			// Convert JS arguments to Go values based on the method's signature
			inputs := make([]reflect.Value, len(args))
			for i := 0; i < len(args); i++ {
				goType := goMethod.Type().In(i)
				jsArg := args[i]
				convertedArg := jsValueToGoType(jsArg, goType)
				if !convertedArg.IsValid() {
					return fmt.Errorf("invalid argument type for argument %d in method %s", i, methodName)
				}
				inputs[i] = convertedArg
			}
			// Call the method with converted arguments
			results := goMethod.Call(inputs)
			// Returning first result (if any) as a JS value
			if len(results) > 0 {
				return goValueToJs(results[0])
			}
			return nil
		})
		// Store the JS function in the result map
		result[lowercaseMethodName] = jsMethod
	}
}

// Helper function to convert JS values to Go types
func jsValueToGoType(jsVal js.Value, goType reflect.Type) reflect.Value {
	switch goType.Kind() {
	case reflect.Int:
		return reflect.ValueOf(jsVal.Int())
	case reflect.Float64:
		return reflect.ValueOf(jsVal.Float())
	case reflect.Bool:
		return reflect.ValueOf(jsVal.Bool())
	case reflect.String:
		return reflect.ValueOf(jsVal.String())
	case reflect.Slice:
		// TODO: Handle slices
		return reflect.ValueOf(jsVal)
	case reflect.Interface:
		return reflect.ValueOf(jsVal)
	}
	return reflect.Value{}
}

// Helper function to convert Go values to JS values
func goValueToJs(goVal reflect.Value) any {
	switch goVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return js.ValueOf(goVal.Int())
	case reflect.Float32, reflect.Float64:
		return js.ValueOf(goVal.Float())
	case reflect.Bool:
		return js.ValueOf(goVal.Bool())
	case reflect.String:
		return js.ValueOf(goVal.String())
	case reflect.Slice:
		return ToJsArray(goVal.Slice(0, goVal.Len()))
	}
	js.Global().Get("console").Call("warn", "Unknown type: "+goVal.Kind().String())
	return js.Null()
}

// Helper function to process slices
func processSlice(v reflect.Value) any {
	var result []any
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Struct {
			result = append(result, marshalToCamelCaseWithMethods(elem, true))
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
func processMap(v reflect.Value) JSON {
	result := make(JSON)
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		if val.Kind() == reflect.Struct {
			result[key.String()] = marshalToCamelCaseWithMethods(val, true)
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
