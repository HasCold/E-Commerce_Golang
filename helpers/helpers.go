package helpers

import (
	"reflect"
)

func IsEmpty(body interface{}) bool {
	val := reflect.ValueOf(body)

	// body := map[string]interface{}{
	// 	"body": "Hello",
	// }

	switch val.Kind() { // val.Kind() tells us the "kind" of value (e.g., reflect.Map for a map, reflect.Slice for a slice, etc.).

	case reflect.Slice, reflect.Chan, reflect.Map, reflect.Array:
		return val.Len() == 0 // The Len() method of reflect.Value gives the number of entries in the map.

	// For pointers or interfaces, it checks whether they are nil. A nil pointer or interface is considered empty.
	case reflect.Interface, reflect.Pointer:
		return val.IsNil()
	}

	return false
}
