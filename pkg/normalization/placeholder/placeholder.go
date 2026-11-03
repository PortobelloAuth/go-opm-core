package placeholder

import (
	"strings"
	"unicode/utf8"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
	"golang.org/x/text/runes"
)

var notAsciiDigit = runes.Predicate(func(r rune) bool {
	return r < '0' || r > '9'
})

// PlaceholderCheck is an interface that defines a placeholder check so
// that each list can be made of different types of check that each
// implement their own logic.
type PlaceholderCheck interface {
	Matches(string) bool
}

type ContextualPlaceholderCheck interface {
	WithContext(normalization.ContextFunc) PlaceholderCheck
}

// StringCheck checks to see if the input string is equal to the value to be
// checked, ignoring case
type StringCheck struct {
	Value string
}

func (sc *StringCheck) Matches(input string) bool {
	return strings.EqualFold(sc.Value, input)
}

// ContainsStringCheck checks to see if the input string contains the value
// to be checked, ignoring case
type ContainsStringCheck struct {
	Value string
}

func (c *ContainsStringCheck) Matches(input string) bool {
	return strings.Contains(strings.ToUpper(input), strings.ToUpper(c.Value))
}

// PrefixStringCheck checks to see if the input string starts with the value
// to be checked, ignoring case
type PrefixStringCheck struct {
	Value string
}

func (p *PrefixStringCheck) Matches(input string) bool {
	return strings.HasPrefix(strings.ToUpper(input), strings.ToUpper(p.Value))
}

// SuffixStringCheck checks to see if the input string ends with the value
// to be checked, ignoring case
type SuffixStringCheck struct {
	Value string
}

func (s *SuffixStringCheck) Matches(input string) bool {
	return strings.HasSuffix(strings.ToUpper(input), strings.ToUpper(s.Value))
}

// RepeatingCharacterCheck checks to see if the input string is all the same
// repeating character
type RepeatingCharacterCheck struct{}

func (r *RepeatingCharacterCheck) Matches(input string) bool {
	if utf8.RuneCountInString(input) < 2 {
		return false
	}
	var first rune
	for i, next := range input {
		if i == 0 {
			first = next
			continue
		}
		if first != next {
			return false
		}
	}
	return true
}

// DigitsCheck checks to see if the input string is entirely made of digits
type DigitsCheck struct{}

func (d *DigitsCheck) Matches(input string) bool {
	if len(input) == 0 {
		return false
	}
	for _, r := range input {
		if notAsciiDigit.Contains(r) {
			return false
		}
	}
	return true
}

// TransformCheck checks to see if the input string matches after it is transformed
type TransformCheck struct {
	Fn    func(string) (string, error)
	Check PlaceholderCheck
	ctx   normalization.ContextFunc
}

func (t *TransformCheck) Matches(input string) bool {
	transformed := input
	check := t.Check
	if check == nil {
		return false
	}
	var err error
	if t.Fn != nil {
		transformed, err = t.Fn(input)
		if err != nil {
			return false
		}
	}
	if t.ctx != nil {
		with, ok := check.(ContextualPlaceholderCheck)
		if ok {
			check = with.WithContext(t.ctx)
		}
	}
	return check.Matches(transformed)
}

func (t *TransformCheck) WithContext(cfn normalization.ContextFunc) PlaceholderCheck {
	with := &TransformCheck{
		Fn:    t.Fn,
		Check: t.Check,
		ctx:   cfn,
	}
	return &FnCheck{
		Fn: with.Matches,
	}
}

// LengthCheck checks to see if the input string is exactly Len characters
type LengthCheck struct {
	Len int
}

func (l *LengthCheck) Matches(input string) bool {
	// make sure to check the character count, not the byte count
	return utf8.RuneCountInString(input) == l.Len
}

// AndCheck checks to see if the input string matches all of multiple checks
type AndCheck struct {
	Checks []PlaceholderCheck
	ctx    normalization.ContextFunc
}

func (a *AndCheck) Matches(input string) bool {
	r := true
	for _, c := range a.Checks {
		if a.ctx != nil {
			with, ok := c.(ContextualPlaceholderCheck)
			if ok {
				c = with.WithContext(a.ctx)
			}
		}
		if !c.Matches(input) {
			return false
		}
	}
	return r
}

func (a *AndCheck) WithContext(cfn normalization.ContextFunc) PlaceholderCheck {
	with := &AndCheck{
		Checks: a.Checks,
		ctx:    cfn,
	}
	return &FnCheck{
		Fn: with.Matches,
	}
}

// OrCheck checks to see if the input string matches any of multiple checks
type OrCheck struct {
	Checks []PlaceholderCheck
	ctx    normalization.ContextFunc
}

func (o *OrCheck) Matches(input string) bool {
	r := false
	for _, c := range o.Checks {
		if o.ctx != nil {
			with, ok := c.(ContextualPlaceholderCheck)
			if ok {
				c = with.WithContext(o.ctx)
			}
		}
		if c.Matches(input) {
			return true
		}
	}
	return r
}

func (o *OrCheck) WithContext(cfn normalization.ContextFunc) PlaceholderCheck {
	with := &OrCheck{
		Checks: o.Checks,
		ctx:    cfn,
	}
	return &FnCheck{
		Fn: with.Matches,
	}
}

// NotCheck checks to see if the input string does not match a check
type NotCheck struct {
	Check PlaceholderCheck
	ctx   normalization.ContextFunc
}

func (n *NotCheck) Matches(input string) bool {
	c := n.Check
	if n.ctx != nil {
		with, ok := c.(ContextualPlaceholderCheck)
		if ok {
			c = with.WithContext(n.ctx)
		}
	}

	return !c.Matches(input)
}

func (n *NotCheck) WithContext(cfn normalization.ContextFunc) PlaceholderCheck {
	with := &NotCheck{
		Check: n.Check,
		ctx:   cfn,
	}
	return &FnCheck{
		Fn: with.Matches,
	}
}

// FnCheck checks to see if a provided function returns true
type FnCheck struct {
	Fn func(string) bool
}

func (fn *FnCheck) Matches(input string) bool {
	if fn.Fn == nil {
		return false
	}
	return fn.Fn(input)
}

// WithContextCheck checks to see if some other aspect of the demographics context matches
type WithContextCheck struct {
	Fn  func(normalization.ContextFunc) bool
	ctx normalization.ContextFunc
}

func (w *WithContextCheck) Matches(input string) bool {
	ctx := w.ctx
	if ctx == nil {
		ctx = func(string) (string, bool) { return "", false }
	}
	return w.Fn(ctx)
}

func (w *WithContextCheck) WithContext(cfn normalization.ContextFunc) PlaceholderCheck {
	with := &WithContextCheck{
		Fn:  w.Fn,
		ctx: cfn,
	}
	return &FnCheck{
		Fn: with.Matches,
	}
}

type Placeholder struct {
	ListLabel string
	List      []PlaceholderCheck
	ctx       normalization.ContextFunc
}

// New allows the caller to create their own Placeholder normalization using a
// supplied label and list. Labels should be unique, primarily to aid in debugging.
func New(label string, list []PlaceholderCheck) *Placeholder {
	return &Placeholder{
		ListLabel: label,
		List:      list,
	}
}

func (p *Placeholder) Apply(input string) (string, error) {
	for _, check := range p.List {
		if p.ctx != nil {
			with, ok := check.(ContextualPlaceholderCheck)
			if ok {
				check = with.WithContext(p.ctx)
			}
		}
		if check.Matches(input) {
			return "", nil
		}
	}
	return input, nil
}

func (p *Placeholder) WithContext(cfn normalization.ContextFunc) normalization.Normalization {
	with := &Placeholder{
		ListLabel: p.ListLabel,
		List:      p.List,
		ctx:       cfn,
	}
	return &normalization.Applicable{Fn: with.Apply}
}
