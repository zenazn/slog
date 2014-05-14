package slog

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func Format(line map[string]interface{}) string {
	keys := make([]string, 0, len(line))
	for k := range line {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		val := fmt.Sprintf("%+v", line[k])

		if strings.IndexFunc(k, unicode.IsSpace) >= 0 {
			k = fmt.Sprintf("%q", k)
		}

		// Reuse storage because I'm a terrible person.
		keys[i] = fmt.Sprintf("%s=%q", k, val)
	}
	return strings.Join(keys, " ") + "\n"
}
