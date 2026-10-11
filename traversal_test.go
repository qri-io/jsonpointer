package jsonpointer

import (
	"encoding/json"
	"errors"
	"testing"
)

type A struct {
	Foo  string `json:"$foo,omitempty"`
	Foo2 bool   `json:"$foo2"`
	Bar  B
}

type B struct {
	nope string
	Baz  int
	C    C
	D    D
}

type C struct {
	Nope string
}

func (c C) JSONProps() map[string]interface{} {
	return map[string]interface{}{
		"bat":   "book",
		"stuff": false,
		"other": nil,
	}
}

type D []string

var data = []byte(`{
  "$foo" : "fooval",
  "$foo2" : true,
  "bar" : {
    "baz" : 1,
    "C" : {
      "won't" : "register"
    }
  }
}`)

func TestWalkJSON(t *testing.T) {
	a := &A{}
	if err := json.Unmarshal(data, a); err != nil {
		t.Errorf("unexpected unmarshal error: %s", err.Error())
		return
	}

	elements := 0
	expectElements := 9
	WalkJSON(a, func(elem interface{}) error {
		t.Logf("%#v", elem)
		elements++
		return nil
	})

	if elements != expectElements {
		t.Errorf("expected %d elements, got: %d", expectElements, elements)
	}

}

func TestWalkJSONStopsOnVisitorError(t *testing.T) {
	values := struct{ First, Second string }{"first", "second"}
	cases := []struct {
		name string
		tree interface{}
	}{
		{"struct", values},
		{"pointer", &values},
		{"map", map[string]string{"first": "value", "second": "value"}},
		{"slice", []string{"first", "second"}},
		{"array", [2]string{"first", "second"}},
		{"json map", map[string]interface{}{"first": "value", "second": "value"}},
		{"json slice", []interface{}{"first", "second"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := errors.New("stop walking")
			calls := 0
			err := WalkJSON(tc.tree, func(elem interface{}) error {
				calls++
				if calls == 2 {
					return want
				}
				return nil
			})
			if err != want {
				t.Errorf("got error %v, want %v", err, want)
			}
			if calls != 2 {
				t.Errorf("visited %d elements, want 2", calls)
			}
		})
	}
}
