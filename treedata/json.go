package treedata

import "fmt"

// SanitizeForJSON will take an interface and make a deep copy of
// it, replacing mam keys with string representations
// this will allow the datastructure to be written as JSON
func SanitizeForJSON(in interface{}) interface{} {
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
		return t

	case []interface{}:
		clone := make([]interface{}, 0)

		for _, value := range t {
			clone = append(clone, SanitizeForJSON(value))
		}
		return clone

	case map[interface{}]interface{}:
		clone := make(map[string]interface{})

		for key, value := range t {
			clone[key.(string)] = SanitizeForJSON(value)
		}
		return clone

	case map[string]interface{}:
		clone := make(map[string]interface{})

		for key, value := range t {
			clone[key] = SanitizeForJSON(value)
		}
		return clone

	default:
		panic(fmt.Sprintf("cannot sanityze %v %T, unsupported type", in, in))
	}
}
