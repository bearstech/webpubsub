package pathmatch

import (
	"testing"
)

func TestPattern(t *testing.T) {
	_, ok := New("a/b/##")
	if ok == nil {
		t.Error()
	}
	_, ok = New("a/b/a#")
	if ok == nil {
		t.Error()
	}
	_, ok = New("a/#/c")
	if ok == nil {
		t.Error()
	}
	_, ok = New("a/#")
	if ok != nil {
		t.Error("This pattern is valid")
	}
	_, ok = New("")
	if ok == nil {
		t.Error()
	}
	_, ok = New("/a/b/c")
	if ok == nil {
		t.Error()
	}
}

func TestPath(t *testing.T) {
	pattern, err := New("a/+/c")
	if err != nil {
		t.Error("Bad pattern : ", err)
	}
	if !pattern.Match("a/b/c") {
		t.Error("Oups")
	}
	if pattern.Match("a/b") {
		t.Error("Oups")
	}
	pattern, _ = New("a/b")
	if !pattern.Match("a/b") {
		t.Error("Oups same stuff")
	}
	pattern, _ = New("a/#")
	if pattern.Match("b") {
		t.Error("Oups same stuff")
	}
	if !pattern.Match("a/b/c") {
		t.Error("Oups same stuff")
	}
}

func TestMultiPath(t *testing.T) {
	pattern, err := New("b/c/d", "a/b/c")
	if err != nil {
		t.Error("Bad pattern : ", err)
	}
	pattern, _ = New("a/b")
	if !pattern.Match("a/b") {
		t.Error("Oups same stuff")
	}
}
