package placeholder

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

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

func (r *DigitsCheck) Matches(input string) bool {
	for _, next := range strings.Split(input, "") {
		if !strings.ContainsAny(next, "0123456789") {
			return false
		}
	}
	return true
}

// LengthCheck checks to see if the input string is exactly Len characters
type LengthCheck struct {
	Len int
}

func (l *LengthCheck) Matches(input string) bool {
	return len(input) == l.Len
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
	fn  func(normalization.ContextFunc) bool
	ctx normalization.ContextFunc
}

func (w *WithContextCheck) Matches(input string) bool {
	if w.ctx == nil {
		w.ctx = func(string) (string, bool) { return "", false }
	}
	return w.fn(w.ctx)
}

func (w *WithContextCheck) WithContext(cfn normalization.ContextFunc) PlaceholderCheck {
	with := &WithContextCheck{
		fn:  w.fn,
		ctx: cfn,
	}
	return &FnCheck{
		Fn: with.Matches,
	}
}

var names = []PlaceholderCheck{
	&StringCheck{"Unknown"},
	&StringCheck{"Unk"},
	&StringCheck{"Test"},
	&StringCheck{"Test Patient"},
	&AndCheck{
		Checks: []PlaceholderCheck{
			&WithContextCheck{
				fn: func(cfn normalization.ContextFunc) bool {
					v, ok := cfn("family")
					return ok && len(v) > 0
				},
			},
			&OrCheck{
				Checks: []PlaceholderCheck{
					&PrefixStringCheck{"Baby Boy"},
					&PrefixStringCheck{"Baby Girl"},
				},
			},
		},
	},
	// &PrefixStringCheck{"Baby Boy"},  // paired with Surname
	// &PrefixStringCheck{"Baby Girl"}, // paired with Surname
	&StringCheck{"Twin A"},
	&StringCheck{"Twin B"},
	&StringCheck{"Newborn"},
	&StringCheck{"N/A"},
	&StringCheck{"None"},
	&StringCheck{"TBD"},
	&AndCheck{
		Checks: []PlaceholderCheck{
			&WithContextCheck{
				fn: func(cfn normalization.ContextFunc) bool {
					v, ok := cfn("family")
					v = strings.ToUpper(v)
					return ok && v == "DOE"
				},
			},
			// paired with generic "John"/"Jane"
			&OrCheck{
				Checks: []PlaceholderCheck{
					&StringCheck{"John"},
					&StringCheck{"Jane"},
				},
			},
		},
	},
	// single-character or repeated-character strings ("X" "ZZZ")
	&RepeatingCharacterCheck{},
}

var dateofbirth = []PlaceholderCheck{
	// NOTE: using string checks here until the rules determine the need for
	// handling various different date formats
	&StringCheck{"01/01/1900"},
	&StringCheck{"01/01/1901"},
	&StringCheck{"01/01/1800"},
	&StringCheck{"01/01/2000"},
	&StringCheck{"11/11/1111"},
	&StringCheck{"01/01/0001"},
	// a DOB equal to the record's registration/admission date (common default-fill error)
}

var last4 = []PlaceholderCheck{
	&NotCheck{Check: &LengthCheck{Len: 4}},
	&StringCheck{"1234"},
	&NotCheck{Check: &DigitsCheck{}},
	&RepeatingCharacterCheck{},
}

var phonenumber = []PlaceholderCheck{
	&StringCheck{"000-000-0000"},
	&StringCheck{"111-111-1111"},
	&StringCheck{"555-555-5555"},
	&StringCheck{"123-456-7890"},
	// any number with a 555 exchange
	&PrefixStringCheck{"555"},
}

var email = []PlaceholderCheck{
	&StringCheck{"test@test.com"},
	// TODO: this needs to remove non-alphanumeric characters first.
	&PrefixStringCheck{"noreply@"}, // include variations with punctuation in the identifier portion of the address
	&StringCheck{"unknown@unknown.com"},
	&StringCheck{"none@none.com"},
	&SuffixStringCheck{"@example.com"},
}

var address = []PlaceholderCheck{
	&StringCheck{"123 Main St"}, // (unaccompanied by a real unit/city match),
	&StringCheck{"PO Box 0"},
	&StringCheck{"Unknown"},
	&StringCheck{"Homeless"},
	&StringCheck{"No Fixed Address"},
	// ZIP 00000 or 99999
	// TODO: use normalization/address/projectusat to parse the address and
	// check the zip (consider if this is too inefficient)
}

var subscriberid = []PlaceholderCheck{
	&StringCheck{"PENDING"},
	&StringCheck{"TBD"},
	&StringCheck{"NONE"},
	// All-zero or all-nine strings
	// FIXME?: this implementation is a little more broad than the rule text
	&RepeatingCharacterCheck{},
	// a payer's documented default/test enrollment value
}

type Placeholder struct {
	List string
	list []PlaceholderCheck
	ctx  normalization.ContextFunc
}

func New(specifier string) (*Placeholder, error) {
	parts := strings.SplitN(specifier, "-", 2)
	if parts[0] != "placeholder" {
		return nil, fmt.Errorf("invalid specifier: %s", specifier)
	}

	p := &Placeholder{}
	switch parts[1] {
	case "names":
		p.list = names
		p.List = parts[1]
	case "dateofbirth":
		p.list = dateofbirth
		p.List = parts[1]
	case "last4":
		p.list = last4
		p.List = parts[1]
	case "phonenumber":
		p.list = phonenumber
		p.List = parts[1]
	case "email":
		p.list = email
		p.List = parts[1]
	case "address":
		p.list = address
		p.List = parts[1]
	case "subscriberid":
		p.list = subscriberid
		p.List = parts[1]
	default:
		// TODO: create logic for building the list from JSON taken from
		// NetworkSettings / RuleSet JSON so that placeholder handling
		// isn't entirely hardcoded.
		return nil, fmt.Errorf("invalid placeholder specifier: %s", specifier)
	}
	return p, nil
}

func (p *Placeholder) Apply(input string) (string, error) {
	for _, check := range p.list {
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
		List: p.List,
		list: p.list,
		ctx:  cfn,
	}
	return &normalization.Applicable{Fn: with.Apply}
}
