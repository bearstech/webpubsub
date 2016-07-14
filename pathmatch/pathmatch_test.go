package pathmatch

import (
	"testing"
)

func TestPattern(t *testing.T) {
	_, ok := NewPattern("a/b/##")
	if ok == nil {
		t.Error()
	}
	_, ok = NewPattern("a/b/a#")
	if ok == nil {
		t.Error()
	}
	_, ok = NewPattern("a/#/c")
	if ok == nil {
		t.Error()
	}
	_, ok = NewPattern("a/#")
	if ok != nil {
		t.Error("This pattern is valid")
	}
	_, ok = NewPattern("")
	if ok == nil {
		t.Error()
	}
	_, ok = NewPattern("/a/b/c")
	if ok == nil {
		t.Error()
	}
}

func TestPath(t *testing.T) {
	pattern, err := NewPattern("a/+/c")
	if err != nil {
		t.Error("Bad pattern : ", err)
	}
	if !pattern.Match("a/b/c") {
		t.Error("Oups")
	}
	if pattern.Match("a/b") {
		t.Error("Oups")
	}
	pattern, _ = NewPattern("a/b")
	if !pattern.Match("a/b") {
		t.Error("Oups same stuff")
	}
}
