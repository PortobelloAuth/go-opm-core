package hashrouted

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/lib/hashing"
	"github.com/PortobelloAuth/go-opm-core/lib/network_settings/hashrouted/hashroutedchallenge"
	"github.com/PortobelloAuth/go-opm-core/lib/rules"
)

type HashRoutedNetworkSettings interface {
	Type() string
	MatchingRules() *rules.RuleSet
	RoutingKeys() []string
	HashFunc() hashing.HashFunc
	RequestCredentialType() string
	AuthenticationStyle() string
	SaltStrategy() string
	AsymmetricEncryptionType() string
	SymmetricEncryptionType() string
}

func New(specifier string) (HashRoutedNetworkSettings, error) {
	residual, found := strings.CutPrefix(specifier, "hash-routed-")
	fmt.Printf("hashrouted.New(%s) with residual: %s\n", specifier, residual)
	if !found {
		return nil, fmt.Errorf("%s is not a hash routed network type", specifier)
	}
	subtype, _, _ := strings.Cut(residual, "-")
	fmt.Printf("  subtype: %s\n", subtype)
	switch subtype {
	case "challenge":
		return hashroutedchallenge.New(specifier)
	}
	return nil, fmt.Errorf("%s is not a supported network type", specifier)
}
