package argparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagSep         string = "|"
	tagKeyValueSep string = "="
)

var validTags = map[string]*regexp.Regexp{
	// "name":  regexp.MustCompile(fmt.Sprintf(`^(name)%s([[:alnum:]-]+)$`, tagKeyValueSep)),
	"name":  regexp.MustCompile(fmt.Sprintf(`^name%s([[:alnum:]-]+)$`, tagKeyValueSep)),
	"pos":   regexp.MustCompile(fmt.Sprintf(`^pos%s(yes)$`, tagKeyValueSep)),
	"help":  regexp.MustCompile(fmt.Sprintf(`^help%s([^%s]+)$`, tagKeyValueSep, tagSep)),
	"nargs": regexp.MustCompile(fmt.Sprintf(`^nargs%s(-?[[:digit:]]+)$`, tagKeyValueSep)),
	// "mutex":      nil,
	// "short":      nil,
}

func parseTags(structTags string) (map[string]string, error) {
	tagValues := make(map[string]string)
	tags := strings.Split(structTags, tagSep)
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
	// 'name' tag must be there
	if tagValues["name"] == "" {
		return nil, fmt.Errorf("name tag is mandatory")
	}
	return tagValues, nil
}

func newArgFromTags(value Value, tags map[string]string) (*Argument, error) {
	var newARg *Argument
	if tags["pos"] == "yes" {
		newARg = NewPosArg(value, tags["help"])
	} else {
		newARg = NewOptArg(value, tags["help"])
	}

	if tags["nargs"] != "" {
		nargs, err := strconv.ParseInt(tags["nargs"], 0, strconv.IntSize)
		if err != nil {
			return nil, formatParseError(tags["nargs"], fmt.Sprintf("%T", int(1)), err)
		}

		err = newARg.SetNArgs(int(nargs))
		if err != nil {
			return nil, err
		}
	}
	return newARg, nil
}
