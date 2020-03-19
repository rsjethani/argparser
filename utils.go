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
	// "nargs":      nil,
	// "short":      nil,
}

func parseTags(structTags string) (map[string]string, error) {
	tagValues := make(map[string]string)
	tags := strings.Split(structTags, "|")
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		invalid := true
		for name, regex := range validTags {
			res := regex.FindStringSubmatch(tag)
			if len(res) == 2 {
				tagValues[name] = res[1]
				invalid = false
			}
		}
		if invalid {
			return nil, fmt.Errorf("invalid tag or value: %s", tag)
		}
	}
	// 'name' tag must be there
	if tagValues["name"] == "" {
		return nil, fmt.Errorf("name tag is mandatory")
	}
	return tagValues, nil
}

func newArgFromTags(value ArgValue, tags map[string]string) (*Argument, error) {
	arg := &Argument{Value: value}
	if tags["pos"] == "yes" {
		arg.Positional = true
	}
	arg.Help = tags["help"]

	if tags["nargs"] != "" {
		val, err := strconv.ParseInt(tags["nargs"], 0, strconv.IntSize)
		if err != nil {
			return nil, err
		}

		err = arg.SetNArgs(int(val))
		if err != nil {
			return nil, err
		}
	}
	return arg, nil
}
