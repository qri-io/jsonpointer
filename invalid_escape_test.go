package jsonpointer

import (
	"reflect"
	"testing"
)

func TestParseTokenEscapes(t *testing.T) {
	for _, text := range []string{"/~", "/~2", "/a~x/b", "/~~0", "#/%7E", "#/~2"} {
		if _, err := Parse(text); err == nil {
			t.Errorf("Parse(%q) succeeded, want invalid escape error", text)
		}
	}
	for _, tc := range []struct {
		text string
		want Pointer
	}{
		{"/~0", Pointer{"~"}}, {"/~1", Pointer{"/"}}, {"/~01", Pointer{"~1"}},
		{"/~0~1", Pointer{"~/"}}, {"/", Pointer{""}}, {"#/a%7E1b", Pointer{"a/b"}},
	} {
		got, err := Parse(tc.text)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Parse(%q)=%v, %v; want %v", tc.text, got, err, tc.want)
		}
	}
	parent := Pointer{"root"}
	if _, err := parent.Descendant("child~2"); err == nil {
		t.Error("Descendant accepted invalid token escape")
	}
	if !reflect.DeepEqual(parent, Pointer{"root"}) {
		t.Fatal("invalid descendant changed parent")
	}
}
