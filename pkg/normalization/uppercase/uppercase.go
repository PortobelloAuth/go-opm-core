package uppercase

import (
	"fmt"
	"strings"
)

type Uppercase struct{}

func New(specifier string) (*Uppercase, error) {
	if specifier != "uppercase" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}
	return &Uppercase{}, nil
}

func (u *Uppercase) Apply(input string) (string, error) {
	result := strings.ToUpper(input)
	return result, nil
}
