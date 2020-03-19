package argparser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var validTags = map[string]*regexp.Regexp{
	"name":  regexp.MustCompile(`^name=([[:alnum:]-]+)$`),
	"pos":   regexp.MustCompile(`^pos=(yes)$`),
	"help":  regexp.MustCompile(`^help=([^|]+)$`),
	"nargs": regexp.MustCompile(`^nargs=(-?[[:digit:]]+)$`),
	// "mutex":      nil,
	// "short":      nil,
}

func parseTags(structTags string) (map[string]string, error) {
	tagValues := make(map[string]string)
	tags := strings.Split(structTags, "|")
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

func newArgFromTags(value ArgValue, tags map[string]string) (*Argument, error) {
	var newARg *Argument
	if tags["pos"] == "yes" {
		newARg = NewPosArg(value, tags["help"])
	} else {
		newARg = NewOptArg(value, tags["help"])
	}

	if tags["nargs"] != "" {
		nargs, err := strconv.ParseInt(tags["nargs"], 0, strconv.IntSize)
		if err != nil {
			return nil, err
		}

		err = newARg.SetNArgs(int(nargs))
		if err != nil {
			return nil, err
		}
	}
	return newARg, nil
}
