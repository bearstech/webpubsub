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
	paths [][]string
}

//NewPattern builds a new pattern and check syntax errors
func New(patterns ...string) (*Pattern, error) {
	px := make([][]string, 0)
	for _, pattern := range patterns {
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
		px = append(px, strings.Split(pattern, "/"))
	}
	return &Pattern{px}, nil
}

func oneMatch(path []string, pattern []string) bool {
	if len(pattern) > len(path) {
		return false
	}
	for i, pat := range path {
		switch pat {
		case "#":
			return true
		case "+":
			continue
		default:
			if i > len(path)-1 {
				return false
			}
			if pat != path[i] {
				return false
			}
		}
	}
	return true

}

func (p *Pattern) Match(path string) bool {
	pp := strings.Split(path, "/")
	for _, path := range p.paths {
		if oneMatch(pp, path) {
			return true
		}
	}
	return false
}
