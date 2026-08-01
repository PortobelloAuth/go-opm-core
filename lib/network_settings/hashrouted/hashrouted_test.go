package hashrouted_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/lib/network_settings/hashrouted"
)

func TestNewHashRoutedNetworkSettings(t *testing.T) {
	unsupported := "not-real-network-type"
	bad, err := hashrouted.New(unsupported)
	if err == nil {
		t.Fatalf("Expected error not found for unsupported network specifier: %s", unsupported)
	}
	if bad != nil {
		t.Fatalf("Expected nil HashRoutedNetworkSettings for unsupported network specifier: %s", unsupported)
	}

	hashroutedchallenge := "hash-routed-challenge"
	hrc, err := hashrouted.New(hashroutedchallenge)
	if err != nil {
		t.Fatalf("Error creating HashRoutedNetworkSettings for network specifier: %s error: %s", hashroutedchallenge, err)
	}
	if hrc == nil {
		t.Fatalf("Failed to create HashRoutedNetworkSettings for network specifier: %s", hashroutedchallenge)
	}
}
