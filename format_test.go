package slog

import "testing"

type foo struct {
	a, b int
}

var formatTests = []struct {
	line map[string]interface{}
	out  string
}{
	{
		map[string]interface{}{"hello": 4, "world": "!"},
		`hello="4" world="!"`,
	},
	{
		map[string]interface{}{"foo": 1.2, "bar": false},
		`bar="false" foo="1.2"`,
	},
	{
		map[string]interface{}{`"`: `"`},
		`"\""="\""`,
	},
	{
		map[string]interface{}{`=`: foo{1, 5}},
		`"="="{a:1 b:5}"`,
	},
	{
		map[string]interface{}{`\`: &foo{2, 3}},
		`"\\"="&{a:2 b:3}"`,
	},
	{
		map[string]interface{}{" hello ": "世界"},
		`" hello "="世界"`,
	},
	{
		map[string]interface{}{"": "hi"},
		`""="hi"`,
	},
}

func TestFormat(t *testing.T) {
	t.Parallel()
	for _, test := range formatTests {
		out := Format(test.line)
		if out != test.out+"\n" {
			t.Errorf("Expected Format(%v) = %q, got %q", test.line,
				test.out, out)
		}
	}
}

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test := formatTests[i%len(formatTests)]
		Format(test.line)
	}
}
