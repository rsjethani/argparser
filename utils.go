package argparser

import (
	"fmt"
	"strings"
)

//if any key is given then its value cannot be empty
//name has to be there
func parseStructTag(tagValue string) (map[string]string, error) {
	tagMap := make(map[string]string)
	for _, value := range strings.Split(tagValue, tagValueSep) {
		parts := strings.Split(value, mapKeyValueSep)
		//TODO: verify key is a proper non-empty string without special symbols etc.
		tagMap[parts[0]] = parts[1]
	}
	if tagMap["name"] == "" {
		return nil, fmt.Errorf("Either 'name' key not specified or its value is empty")
	}
	return tagMap, nil
}
