package plaintextcomponent

import (
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/match"
	"github.com/PortobelloAuth/go-opm-core/pkg/match/matcher/componentmatcher"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

var noopHash hashing.Hash

func init() {
	var err error
	noopHash, err = hashing.New("Noop")
	if err != nil {
		panic(err)
	}
}

type Matcher struct {
	internal *componentmatcher.Matcher
}

func New(rs *rules.RuleSet) (*Matcher, error) {
	c, err := componentmatcher.New(
		rs,
		noopHash.HashFunc,
		// empty salt does not matter for noop
		[]byte{},
		PlaintextHashToString,
	)
	if err != nil {
		return nil, err
	}

	return &Matcher{
		internal: c,
	}, nil
}

// HashDemographics creates HashedDemographics from Demographics by applying the Matcher's RuleSet
// Normalizations. Only demographics identified in the RuleSet are set in the resulting HashedDemographics
// Hashes map.
func (m *Matcher) HashDemographics(raw *match.Demographics) (*match.HashedDemographics, error, error) {
	return m.internal.HashDemographics(raw)
}

func (m *Matcher) Match(calc, comp *match.HashedDemographics) (*match.MatchResult, error) {
	return m.internal.Match(calc, comp)
}

func PlaintextHashToString(b []byte) (string, error) {
	return string(b), nil
}
