package jsutils

import (
	"reflect"
	"syscall/js"
)

type JSON = map[string]interface{}

func ToJsArray(value interface{}) js.Value {
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
