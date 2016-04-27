package webform

import (
	"reflect"
	"testing"
)

func TestNested(t *testing.T) {
	m := map[string][]string{
		"Text":    []string{"one"},
		"Integer": []string{"1"},
		"float":   []string{"1.1"},

		"TextNested":  []string{"two"},
		"IntNested":   []string{"2"},
		"floatNested": []string{"2.2"},

		"TextNested2":  []string{"three"},
		"IntNested2":   []string{"3"},
		"floatNested2": []string{"3.3"},
	}

	type Nested2 struct {
		TextNested2  string
		IntNested2   int
		FloatNested2 float64 `webform:"floatNested2"`
	}

	type Nested struct {
		Nested2
		TextNested  string
		IntNested   int
		FloatNested float64 `webform:"floatNested"`
	}

	type Tp struct {
		Text    string
		Integer int
		Float   float64 `webform:"float"`
		Nested
	}

	var tp Tp
	Decode(&tp, m)
	expected := Tp{"one", 1, 1.1, Nested{Nested2{"three", 3, 3.3}, "two", 2, 2.2}}

	if !reflect.DeepEqual(&tp, &expected) {
		t.Errorf("not equal %v and %v", tp, expected)
	}
}

func TestNestedPtr(t *testing.T) {
	m := map[string][]string{
		"Text":    []string{"one"},
		"Integer": []string{"1"},
		"float":   []string{"1.1"},

		"TextNested":  []string{"two"},
		"IntNested":   []string{"2"},
		"floatNested": []string{"2.2"},

		"TextNested2":  []string{"three"},
		"IntNested2":   []string{"3"},
		"floatNested2": []string{"3.3"},
	}

	type Nested2 struct {
		TextNested2  string
		IntNested2   int
		FloatNested2 float64 `webform:"floatNested2"`
	}

	type Nested struct {
		*Nested2
		TextNested  string
		IntNested   int
		FloatNested float64 `webform:"floatNested"`
	}

	type Tp struct {
		Text    string
		Integer int
		Float   float64 `webform:"float"`
		*Nested
	}

	var tp Tp
	Decode(&tp, m)
	expected := Tp{"one", 1, 1.1, &Nested{&Nested2{"three", 3, 3.3}, "two", 2, 2.2}}

	if !reflect.DeepEqual(&tp, &expected) {
		t.Errorf("not equal %v and %v", tp, expected)
	}
}
