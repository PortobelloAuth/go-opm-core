package hashroutedchallenge

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
	"github.com/PortobelloAuth/go-opm-core/pkg/network_settings"
	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

func init() {
	err := network_settings.Register(
		"hash-routed-challenge",
		func(specifier string) (network_settings.NetworkSettings, error) { return New(specifier) },
	)
	if err != nil {
		panic(err)
	}
}

type HashRoutedChallengeSettings struct {
	// should always be "hash-routed-challenge"
	NetworkType string
	// the set of rules that apply to this network
	MatchRules *rules.RuleSet
	// the demographics that are used to generate routing keys
	RoutingKeyDemographics []string
	// the hash used for routing keys and challenges
	HashFn hashing.HashFunc
	// the types of credentials that the requestor will eventually send
	// to the responder
	RequestCredentials string
	// how requestors, responders, and facilitators authenticate to this network
	NetAuthenticationStyle string
	// how routing key salts are generated
	RoutingSaltStrategy string
	// The type of public keys the network uses
	AsymEncryptionType string
	// The type of symetric encryption the network uses
	SymEncryptionType string
}

func New(specifier string) (*HashRoutedChallengeSettings, error) {
	residual, found := strings.CutPrefix(specifier, "hash-routed-challenge")
	if !found {
		return nil, fmt.Errorf("%s is not a hash routed challenge network type", specifier)
	}
	subtype, _, _ := strings.Cut(residual, "-")
	if len(subtype) == 0 {
		return &HashRoutedChallengeSettings{
			NetworkType: specifier,
		}, nil
	}
	return nil, fmt.Errorf("%s is not a supported hash routed challenge network type", specifier)

}

func (hrc *HashRoutedChallengeSettings) Type() string {
	return hrc.NetworkType
}

func (hrc *HashRoutedChallengeSettings) MatchingRules() *rules.RuleSet {
	return nil
}

func (hrc *HashRoutedChallengeSettings) RoutingKeys() []string {
	return nil
}

func (hrc *HashRoutedChallengeSettings) HashFunc() hashing.HashFunc {
	return nil
}

func (hrc *HashRoutedChallengeSettings) RequestCredentialType() string {
	return ""
}

func (hrc *HashRoutedChallengeSettings) AuthenticationStyle() string {
	return ""
}

func (hrc *HashRoutedChallengeSettings) SaltStrategy() string {
	return ""
}

func (hrc *HashRoutedChallengeSettings) AsymmetricEncryptionType() string {
	return ""
}

func (hrc *HashRoutedChallengeSettings) SymmetricEncryptionType() string {
	return ""
}
