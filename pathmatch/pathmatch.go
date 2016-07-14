/*
Package pathmatch implement MQTT pathmatch rules

http://mosquitto.org/man/mqtt-7.html



*/
package pathmatch

import (
	"errors"
	"strings"
)

// A path pattern
type Pattern struct {
	path []string
}

//NewPattern builds a new pattern and check syntax errors
func NewPattern(pattern string) (*Pattern, error) {
	if len(pattern) == 0 {
		return nil, errors.New("Empty pattern is not valid, use \"#\"")
	}
	if pattern[0] == '/' {
		return nil, errors.New("Leading / is not valid")
	}
	hash := strings.Index(pattern, "#")
	if hash != -1 {
		if hash+1 != len(pattern) {
			return nil, errors.New("# is not ending the pattern")
		}
		if pattern[len(pattern)-2] != '/' {
			return nil, errors.New("# should be the complete name")
		}
	}
	return &Pattern{strings.Split(pattern, "/")}, nil
}

func (p *Pattern) Match(path string) bool {
	pp := strings.Split(path, "/")
	for i, pat := range p.path {
		switch pat {
		case "#":
			return true
		case "+":
			continue
		default:
			if i > len(pp)-1 {
				return false
			}
			if pat != pp[i] {
				return false
			}
		}
	}
	return true
}
