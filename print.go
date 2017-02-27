package gutil

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
)

func PrintAsYAML(obj interface{}) {
	out, err := yaml.Marshal(obj)
	if err != nil {
		log.Printf("Failed to print %T as YAML, %v", err)
	}
	fmt.Print(string(out))
}
