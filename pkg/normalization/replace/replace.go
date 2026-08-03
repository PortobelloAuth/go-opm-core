package replace

import (
	"fmt"
	"strings"
)

type Replace struct {
	From string
	To   string
}

func New(specifier string) (*Replace, error) {
	parts := strings.SplitN(specifier, "/", 4)
	if parts[0] != "s" || len(parts[1]) == 0 || len(parts[2]) == 0 || len(parts[3]) != 0 {
		return nil, fmt.Errorf("Invalid specifier string provided: %s", specifier)
	}
	return &Replace{
		From: parts[1],
		To:   parts[2],
	}, nil
}

func (r *Replace) Apply(input string) (string, error) {
	result := strings.ReplaceAll(input, r.From, r.To)
	return result, nil
}
