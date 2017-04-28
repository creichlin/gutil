package gutil

import "fmt"

// ConvertToJSONTree will take a value or tree of lists and maps
// and return a version where types that are not json compatible
// are converted to json formats.
// Those are:
// int, int32, int64, float32 -> float64
// maps with non string keys to maps with string keys
func ConvertToJSONTree(in interface{}) interface{} {

	if in == nil {
		return nil
	}

	switch t := in.(type) {
	case string:
		return t

	case bool:
		return t

	case float64:
		return t

	case int:
		return float64(t)

	case int32:
		return float64(t)

	case int64:
		return float64(t)

	case float32:
		return float64(t)

	case []interface{}:
		clone := make([]interface{}, 0)

		for _, value := range t {
			clone = append(clone, ConvertToJSONTree(value))
		}
		return clone

	case map[interface{}]interface{}:
		clone := make(map[string]interface{})

		for key, value := range t {
			clone[fmt.Sprint(key)] = ConvertToJSONTree(value)
		}
		return clone

	case map[string]interface{}:
		clone := make(map[string]interface{})

		for key, value := range t {
			clone[key] = ConvertToJSONTree(value)
		}
		return clone

	default:
		panic(fmt.Sprintf("cannot sanityze %v %T, unsupported type", in, in))
	}
}
