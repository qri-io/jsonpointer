package jsonpointer

import "testing"

func TestDescendantDoesNotOverwriteSibling(t *testing.T) {
	for _, parent := range []Pointer{NewPointer(), NewPointer().RawDescendant("parent")} {
		t.Run(parent.String(), func(t *testing.T) {
			prefix := parent.String()
			first, err := parent.Descendant("first")
			if err != nil {
				t.Fatal(err)
			}
			second, err := parent.Descendant("second")
			if err != nil {
				t.Fatal(err)
			}
			if got := first.String(); got != prefix+"/first" {
				t.Errorf("creating sibling changed first descendant to %q", got)
			}
			if got := second.String(); got != prefix+"/second" {
				t.Errorf("second descendant = %q", got)
			}
			if len(parent) > 0 {
				second[0] = "changed"
				if got := parent.String(); got != prefix {
					t.Errorf("editing descendant changed parent to %q", got)
				}
			}
		})
	}
}
