package argparser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplitKV(t *testing.T) {
	data := make(map[string][]string)
	data[""] = []string{}
	data[fmt.Sprintf("%[1]c%[1]c%[1]c", tagSep)] = []string{}
	data[fmt.Sprintf("a%[1]cb%[1]cc%[1]cd", tagSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca%[1]cb%[1]cc%[1]cd%[1]c", tagSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca\\%[1]cb%[1]cc%[1]cd\\%[1]c", tagSep)] = []string{fmt.Sprintf("a%cb", tagSep), "c", fmt.Sprintf("d%c", tagSep)}
	for input, expected := range data {
		if got := splitKV(input, tagSep); !reflect.DeepEqual(got, expected) {
			t.Errorf("testing: splitKV(%q, %c); expected: %q; Got: %q", input, tagSep, expected, got)
		}
	}
}
