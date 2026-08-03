package network_settings_test

import (
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/network_settings"
)

func TestNewNetworkSettings(t *testing.T) {
	unsupported := "not-real-network-type"
	bad, err := network_settings.New(unsupported)
	if err == nil {
		t.Fatalf("Expected error not found for unsupported network specifier: %s", unsupported)
	}
	if bad != nil {
		t.Fatalf("Expected nil NetworkSettings for unsupported network specifier: %s", unsupported)
	}

	hashroutedchallenge := "hash-routed-challenge"
	hrc, err := network_settings.New(hashroutedchallenge)
	if err != nil {
		t.Fatalf("Error creating NetworkSettings for network specifier: %s error: %s", hashroutedchallenge, err)
	}
	if hrc == nil {
		t.Fatalf("Failed to create NetworkSettings for network specifier: %s", hashroutedchallenge)
	}
}
