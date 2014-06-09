package slog

import (
	"reflect"
	"testing"
)

type fakeTime struct{}

func (_ fakeTime) String() string {
	return "now"
}

func TestBasic(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan map[string]interface{}, 1)
	root.SlogTo(target)

	root.Log(Data{"hello": "world"})

	actual := <-target
	expected := map[string]interface{}{
		"$level": LInfo,
		"$time":  fakeTime{},
		"hello":  "world",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %#v, but got %#v", expected, actual)
	}
}

func TestLogTo(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan string, 1)
	root.LogTo(target)

	root.Log(Data{"hello": "world"})

	actual := <-target
	expected := `$level="INFO" $time="now" hello="world"` + "\n"
	if expected != actual {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}
}

func TestMulti(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan string, 100)
	root.LogTo(target)

	go func() {
		for i := 0; i < 100; i++ {
			root.Log(Data{"hello": "world"})
		}
	}()

	expected := `$level="INFO" $time="now" hello="world"` + "\n"
	for i := 0; i < 100; i++ {
		actual := <-target
		if expected != actual {
			t.Errorf("Expected %q, but got %q", expected, actual)
		}
	}
}

func expectLines(t *testing.T, ch <-chan string, lines []string) {
	for _, line := range lines {
		actual := <-ch
		if actual != line {
			t.Errorf("Expected %q, but got %q", line, actual)
		}
	}
}

func TestLevels(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan string, 4)
	root.LogTo(target)

	root.Debug(Data{})
	root.Log(Data{})
	root.Warn(Data{})
	root.Error(Data{})

	expectLines(t, target, []string{
		`$level="INFO" $time="now"` + "\n",
		`$level="WARN" $time="now"` + "\n",
		`$level="ERROR" $time="now"` + "\n",
	})
}

func TestRules(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan string, 4)
	root.LogTo(target)
	root.SetLevel("example.com/not/a/real/thing", LError)
	root.SetLevel("github.com/zenazn/slog", LWarn)

	root.Debug(Data{})
	root.Log(Data{})
	root.Warn(Data{})
	root.Error(Data{})

	expectLines(t, target, []string{
		`$level="WARN" $time="now"` + "\n",
		`$level="ERROR" $time="now"` + "\n",
	})

	root.SetLevel("github.com/zenazn/slog.", LDebug)

	root.Debug(Data{})
	root.Log(Data{})
	root.Warn(Data{})
	root.Error(Data{})

	expectLines(t, target, []string{
		`$level="DEBUG" $time="now"` + "\n",
		`$level="INFO" $time="now"` + "\n",
		`$level="WARN" $time="now"` + "\n",
		`$level="ERROR" $time="now"` + "\n",
	})
}

func TestBasicBind(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan string, 2)
	root.LogTo(target)

	sub := root.Bind(Data{"hello": "world"})
	sub.Log(Data{"foo": "bar"})
	sub.Log(Data{"hello": "universe", "color": "red"})

	expectLines(t, target, []string{
		`$level="INFO" $time="now" foo="bar" hello="world"` + "\n",
		`$level="INFO" $time="now" color="red" hello="universe"` + "\n",
	})
}

func TestBind(t *testing.T) {
	t.Parallel()

	root := makeRoot(fakeTime{})
	target := make(chan string, 7)
	target2 := make(chan string, 3)
	root.LogTo(target)

	sub := root.Bind(Data{"hello": "world"})
	subsub := sub.Bind(Data{"foo": "bar"})

	root.Log(Data{"space": "ship"})
	sub.Log(Data{"space": "ship"})
	subsub.Log(Data{"space": "ship"})

	sub.LogTo(target2, LInfo)
	root.Log(Data{"space": "ship"})
	sub.Log(Data{"space": "ship"})
	subsub.Log(Data{"space": "ship"})

	sub.SetLevel("github.com", LWarn)
	root.Log(Data{"space": "ship"})
	sub.Log(Data{"space": "ship"})
	subsub.Log(Data{"space": "ship"})

	root.SetLevel("github.com/zenazn/slog", LDebug)
	subsub.LogTo(target, LInfo)
	root.Log(Data{"space": "ship"})
	subsub.Log(Data{"space": "ship"})
	sub.Log(Data{"space": "ship"})

	prefix := `$level="INFO" $time="now" `
	expectLines(t, target, []string{
		prefix + `space="ship"` + "\n",
		prefix + `hello="world" space="ship"` + "\n",
		prefix + `foo="bar" hello="world" space="ship"` + "\n",
		prefix + `space="ship"` + "\n",
		prefix + `space="ship"` + "\n",
		prefix + `space="ship"` + "\n",
		prefix + `foo="bar" hello="world" space="ship"` + "\n",
	})
	expectLines(t, target2, []string{
		prefix + `hello="world" space="ship"` + "\n",
		prefix + `foo="bar" hello="world" space="ship"` + "\n",
		prefix + `hello="world" space="ship"` + "\n",
	})
}
