package controllers

import (
	"strconv"
)

// ParseUintFromPath extracts the trailing integer id from a path like /resource/123
func ParseUintFromPath(path string) (uint, error) {
	idx := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			idx = i
			break
		}
	}
	if idx == -1 || idx == len(path)-1 {
		return 0, strconv.ErrSyntax
	}
	v, err := strconv.Atoi(path[idx+1:])
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
