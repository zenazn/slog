package slog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func needsQuote(r rune) bool {
	return r == ' ' || r == '"' || r == '\\' || r == '=' ||
		!unicode.IsPrint(r)
}

// Format formats a generic map into string form. It sorts all key-value pairs
// by increasing key, and prints strings of the form "key1=value1 key2=value2"
// etc. The values are first formatted using package fmt's "%+v" encoding, then
// quoted using double quotes. The keys are quoted if they'd be ambiguous.
func Format(data map[string]interface{}) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		val := fmt.Sprintf("%+v", data[k])

		if strings.IndexFunc(k, needsQuote) >= 0 {
			k = strconv.Quote(k)
		}

		// Reuse storage because I'm a terrible person.
		keys[i] = k + "=" + strconv.Quote(val)
	}
	return strings.Join(keys, " ") + "\n"
}
