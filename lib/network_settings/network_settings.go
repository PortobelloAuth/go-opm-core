package network_settings

import (
	"github.com/PortobelloAuth/go-opm-core/lib/hashing"
)

type NetworkSettings interface {
	Type() string
	MatchingRules() []string
	RoutingKeys() []string
	HashFunc() hashing.HashFunc
	RequestCredentialType() string
	AuthenticationStyle() string
	SaltStrategy() string
	AsymmetricEncryptionType() string
	SymmetricEncryptionType() string
}

type HashRoutedChallengeSettings struct {
	Type                     string
	MatchingRules            []string
	RoutingKeys              []string
	HashFunc                 hashing.HashFunc
	RequestCredentialType    string
	AuthenticationStyle      string
	SaltStrategy             string
	AsymmetricEncryptionType string
	SymmetricEncryptionType  string
}
