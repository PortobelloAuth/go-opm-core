package normalization

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/digits"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/dob"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/lowercase"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/phone"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/punctuation"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/replace"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/transliterate"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/uppercase"
	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/whitespace"
)

type Normalization interface {
	Apply(input string) (string, error)
}

type NormalizationConstructor func(specifier string) (Normalization, error)

var normalizationRegistery map[string]NormalizationConstructor
var newCalled = false

func init() {
	normalizationRegistery = make(map[string]NormalizationConstructor)

	// register the built-in Normalizations
	builtins := []struct {
		K string
		N func(specifier string) (Normalization, error)
	}{
		{"digits_only", func(specifier string) (Normalization, error) { return digits.New(specifier) }},
		{"dob_iso", func(specifier string) (Normalization, error) { return dob.New(specifier) }},
		{"lowercase", func(specifier string) (Normalization, error) { return lowercase.New(specifier) }},
		{"phone_e164", func(specifier string) (Normalization, error) { return phone.New(specifier) }},
		{"punctuation", func(specifier string) (Normalization, error) { return punctuation.New(specifier) }},
		{"transliterate", func(specifier string) (Normalization, error) { return transliterate.New(specifier) }},
		{"uppercase", func(specifier string) (Normalization, error) { return uppercase.New(specifier) }},
		{"us-101-key", func(specifier string) (Normalization, error) {
			return NewSequence([]string{
				"transliterate",
				"lowercase",
			})
		}},
		{"us-reduced_v0", func(specifier string) (Normalization, error) {
			return NewSequence([]string{
				"transliterate",
				"punctuation",
				"whitespace",
				"lowercase",
			})
		}},
		{"whitespace", func(specifier string) (Normalization, error) { return whitespace.New(specifier) }},
	}

	for _, b := range builtins {
		normalizationRegistery[b.K] = b.N
	}

	prefixBuiltins := []struct {
		K string
		N func(specifier string) (Normalization, error)
	}{
		{"s", func(specifier string) (Normalization, error) { return replace.New(specifier) }},
	}

	for _, p := range prefixBuiltins {
		k := prefixKey(p.K)
		normalizationRegistery[k] = p.N
	}
}

func prefixKey(k string) string {
	return fmt.Sprintf("prefix/%s", k)
}

func Register(specifier string, fn NormalizationConstructor, isPrefix bool) error {
	// FIXME?: don't allow this to be called after New()? It makes testing annoying.
	if newCalled {
		return fmt.Errorf("Register() called after New()")
	}

	k := specifier
	if isPrefix {
		k = prefixKey(specifier)
	}

	// FIXME?: don't allow overwriting builtins?
	normalizationRegistery[k] = fn

	return nil
}

// TODO: consider removing the specifier as an argument for normalizations that don't need it
// TODO: consider whether we need short prefixes for all normalization specifiers
// TODO: figure out an interaction with rules that scores changes made by a normalization (or
// treat changes that should be "scored" like replace as a distinct concept: transformations)

// New creates a new Normalization according to the specifier string
func New(specifier string) (Normalization, error) {
	newCalled = true
	prefix, _, found := strings.Cut(specifier, "/")
	k := specifier
	if found {
		k = prefixKey(prefix)
	}
	fn, ok := normalizationRegistery[k]
	if ok {
		return fn(specifier)
	}

	return nil, fmt.Errorf("%s is not a registered normalization", specifier)
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

	// guard against rule.Rule.Normalizations not being initialized when
	// when there are no normalizations at that level.
	if seq != nil {
		var err error
		for _, norm := range seq.Seq {
			if norm == nil {
				// in case we haven't filtered normalizations effectively
				// when building the sequence
				continue
			}
			result, err = norm.Apply(result)
			if err != nil {
				return "", err
			}
		}
	}
	return result, nil
}
