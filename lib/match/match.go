package match

import (
	"github.com/PortobelloAuth/go-opm-core/lib/hashing"
)

type MatchResult struct {
	Matched   bool
	Threshold float32
	// Confidence Interval / Error
}

type MatchSettings struct {
	HashAlgorithm hashing.HashFunc
}

type DemographicSet struct {
	Settings MatchSettings
}

func Match() *MatchResult {
	// Failed to match with no confidence
	return &MatchResult{
		Matched: false,
	}
}
