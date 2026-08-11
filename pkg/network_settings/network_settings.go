package network_settings

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/rules"
)

type NetworkSettingsConstructor func(specifier string) (NetworkSettings, error)

var netSettingsRegistery map[string]NetworkSettingsConstructor
var newCalled = false

func init() {
	// NOTE: NetworkSettings implementations should register themselves
	netSettingsRegistery = make(map[string]NetworkSettingsConstructor)
}

func Register(specifier string, fn NetworkSettingsConstructor) error {
	// FIXME?: don't allow this to be called after New()? It makes testing annoying.
	if newCalled {
		return fmt.Errorf("Register() called after New()")
	}

	// FIXME?: don't allow overwriting builtins? If not, how would a package user test out replacement code?
	spec := strings.ToLower(specifier)
	// _, found := netSettingsRegistery[spec]
	// if found {
	// 	return fmt.Errorf("Registration already exists for %s", spec)
	// }
	netSettingsRegistery[spec] = fn

	return nil
}

type NetworkSettings interface {
	Type() string
	MatchingRules() *rules.RuleSet
	RequestCredentialType() string
	AuthenticationStyle() string
}

func New(specifier string) (NetworkSettings, error) {
	newCalled = true

	sep := "-"
	parts := strings.Split(strings.ToLower(specifier), sep)

	for len(parts) > 0 {
		prefix := strings.Join(parts, sep)
		fn, ok := netSettingsRegistery[prefix]
		if ok {
			return fn(specifier)
		}
		parts = parts[0 : len(parts)-1]
	}

	return nil, fmt.Errorf("%s is not a supported network type", specifier)
}

// Go's UnmarshalJSON API can't unmarshall to an interface directly so we create a wrapper type
// that is used for unmarshalling any kind of registered NetworkSettings in to its Settings field
type TypedNetworkSettings struct {
	Type     string
	Settings *NetworkSettings `json:"-"` // Ignore standard unmarshaling for everything except "type"
}

func (u *TypedNetworkSettings) UnmarshalJSON(data []byte) error {
	type Alias TypedNetworkSettings
	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return fmt.Errorf("Error parsing TypedNetworkSettings: %w", err)
	}

	u.Type = a.Type

	// Call New() to get a NetworkSettings instance of the specified type
	n, err := New(u.Type)
	if err != nil {
		return fmt.Errorf("Error instantiating NetworkSettings: %w", err)
	}

	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("Error parsing NetworkSettings: %w", err)
	}

	return nil
}
