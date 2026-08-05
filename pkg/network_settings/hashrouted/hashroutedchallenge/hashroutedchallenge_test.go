package hashroutedchallenge_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/network_settings/hashrouted/hashroutedchallenge"
)

func TestNewHashRoutedChallenge(t *testing.T) {
	unsupported := "not-real-network-type"
	bad, err := hashroutedchallenge.New(unsupported)
	if err == nil {
		t.Fatalf("Expected error not found for unsupported network specifier: %s", unsupported)
	}
	if bad != nil {
		t.Fatalf("Expected nil HashRoutedNetworkSettings for unsupported network specifier: %s", unsupported)
	}

	hrctype := "hash-routed-challenge"
	hrc, err := hashroutedchallenge.New(hrctype)
	if err != nil {
		t.Fatalf("Error creating HashRoutedNetworkSettings for network specifier: %s error: %s", hrctype, err)
	}
	if hrc == nil {
		t.Fatalf("Failed to create HashRoutedNetworkSettings for network specifier: %s", hrctype)
	}
}
