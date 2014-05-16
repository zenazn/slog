package slog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func needsQuote(r rune) bool {
	return r == ' ' || r == '"' || r == '\\' || !unicode.IsPrint(r)
}

func Format(line map[string]interface{}) string {
	keys := make([]string, 0, len(line))
	for k := range line {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		val := fmt.Sprintf("%+v", line[k])

		if strings.IndexFunc(k, needsQuote) >= 0 {
			k = strconv.Quote(k)
		}

		// Reuse storage because I'm a terrible person.
		keys[i] = k + "=" + strconv.Quote(val)
	}
	return strings.Join(keys, " ") + "\n"
}
