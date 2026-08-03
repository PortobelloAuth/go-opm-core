package normalization

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/lowercase"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/punctuation"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/replace"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/transliterate"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/uppercase"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/whitespace"
)

type Normalization interface {
	Apply(input string) (string, error)
}

// TODO: consider removing the specifier as an argument for normalizations that don't need it
// TODO: consider whether we need short prefixes for all normalization specifiers
// TODO: figure out an interaction with rules that scores changes made by a normalization (or
// treat changes that should be "scored" like replace as a distinct concept: transformations)

// New creates a new Normalization according to the specifier string
func New(specifier string) (Normalization, error) {
	switch specifier {
	case "lowercase":
		return lowercase.New(specifier)
	case "punctuation":
		return punctuation.New(specifier)
	case "transliterate":
		return transliterate.New(specifier)
	case "uppercase":
		return uppercase.New(specifier)
	case "us-101-key":
		return NewSequence([]string{"transliterate", "lowercase"})
	case "us-reduced_v0":
		return NewSequence([]string{"transliterate", "punctuation", "whitespace", "lowercase"})
	case "whitespace":
		return whitespace.New(specifier)
	default:
		prefix, _, _ := strings.Cut(specifier, "/")
		switch prefix {
		case "s":
			return replace.New(specifier)
		}

		return nil, fmt.Errorf("%s is not a supported normalization", specifier)
	}
}

type NormalizationSequence struct {
	Seq []Normalization
}

func NewSequence(order []string) (*NormalizationSequence, error) {
	sequence := make([]Normalization, 0)
	for _, s := range order {
		norm, err := New(s)
		if err != nil {
			return nil, err
		}
		if norm == nil {
			return nil, fmt.Errorf("unable to initialize normalization for %s", s)
		}
		sequence = append(sequence, norm)
	}

	return &NormalizationSequence{
		Seq: sequence,
	}, nil
}

func (seq *NormalizationSequence) Apply(input string) (string, error) {
	result := input
	var err error
	for _, norm := range seq.Seq {
		result, err = norm.Apply(result)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}
