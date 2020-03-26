package argparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagSep         rune = ','
	tagKeyValueSep rune = '='
)

var validTags = map[string]*regexp.Regexp{
	// "name":  regexp.MustCompile(fmt.Sprintf(`^(name)%s([[:alnum:]-]+)$`, tagKeyValueSep)),
	"name":  regexp.MustCompile(fmt.Sprintf(`^name%c([[:alnum:]-]+)$`, tagKeyValueSep)),
	"type":  regexp.MustCompile(fmt.Sprintf(`^type%c(pos|opt|switch)$`, tagKeyValueSep)),
	"help":  regexp.MustCompile(fmt.Sprintf(`^help%c(.+)$`, tagKeyValueSep)),
	"nargs": regexp.MustCompile(fmt.Sprintf(`^nargs%c(-?[[:digit:]]+)$`, tagKeyValueSep)),
	// "mutex":      nil,
	// "short":      nil,
}

func splitKV(src string, sep rune) []string {
	backSlash := '\\'
	parts := make([]string, 0)
	b := &strings.Builder{}
	var prevRune rune
	for _, curRune := range src {
		switch {
		// rune is a backSlash simply skip it
		case curRune == backSlash:
		// rune is a sep but it is not escaped by backSlash
		case curRune == sep && prevRune != backSlash:
			if b.Len() != 0 {
				parts = append(parts, b.String())
				b.Reset()
			}
		// rune is either not a sep/backslash or if it is a sep then it is escaped by backskash
		default:
			b.WriteRune(curRune)
		}
		prevRune = curRune
	}
	// append any remaining runes between last sep and end of src
	if b.Len() != 0 {
		parts = append(parts, b.String())
	}
	return parts
}

func parseTags(structTags string) (map[string]string, error) {
	tagValues := make(map[string]string)
	tags := splitKV(structTags, tagSep)
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		unknownTag := true
		for name, regex := range validTags {
			res := regex.FindStringSubmatch(tag)
			if len(res) == 2 {
				tagValues[name] = res[1]
				unknownTag = false
			}
		}
		if unknownTag {
			return nil, fmt.Errorf("unknown tag and/or invalid value: %s", tag)
		}
	}

	return tagValues, nil
}

func newArgFromTags(value Value, fieldName string, structTags string) (*Argument, string, error) {
	tags, err := parseTags(structTags)
	if err != nil {
		return nil, "", err
	}

	var newARg *Argument
	// calculate nargs
	switch tags["type"] {
	case "pos":
		if tags["nargs"] == "" {
			tags["nargs"] = "1"
		}
		newARg = NewPosArg(value, tags["help"])
	case "switch":
		if tags["nargs"] == "" {
			tags["nargs"] = "0"
		} else if tags["nargs"] != "0" {
			return nil, "", fmt.Errorf("nargs can only be 0 for type=switch")
		}
		newARg = NewOptArg(value, tags["help"])
	case "opt", "":
		if tags["nargs"] == "" {
			tags["nargs"] = "1"
		}
		newARg = NewOptArg(value, tags["help"])
	}
	nargs, err := strconv.ParseInt(tags["nargs"], 0, strconv.IntSize)
	if err != nil {
		return nil, "", formatParseError(tags["nargs"], fmt.Sprintf("%T", int(1)), err)
	}

	err = newARg.SetNArgs(int(nargs))
	if err != nil {
		return nil, "", err
	}

	// calculate name: if "name" not specified then simlpy use field's name in lower case
	if tags["name"] == "" {
		tags["name"] = strings.ToLower(fieldName)
	}

	return newARg, tags["name"], nil
}
