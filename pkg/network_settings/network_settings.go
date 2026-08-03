package network_settings

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/network_settings/hashrouted"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

type NetworkSettings interface {
	Type() string
	MatchingRules() *rules.RuleSet
	RequestCredentialType() string
	AuthenticationStyle() string
}

func New(specifier string) (NetworkSettings, error) {
	networktype, remainder, _ := strings.Cut(specifier, "-")
	switch networktype {
	case "hash":
		hashroutingtype, _, _ := strings.Cut(remainder, "-")
		switch hashroutingtype {
		case "routed":
			return hashrouted.New(specifier)
		}
	}
	return nil, fmt.Errorf("%s is not a supported network type", specifier)
}
