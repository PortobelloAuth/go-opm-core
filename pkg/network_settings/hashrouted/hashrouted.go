package hashrouted

import (
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
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
