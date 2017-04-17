package gutil

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

// PrintAsYAML prins the given object formatted as yaml to stdout
func PrintAsYAML(obj interface{}) {
	out, err := yaml.Marshal(obj)
	if err != nil {
		log.Printf("Failed to print %T as YAML, %v", obj, err)
	}
	fmt.Println(string(out))
}
