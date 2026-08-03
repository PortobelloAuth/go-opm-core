package challenges

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/challenge/x2r"
)

type Challenge interface {
	ChallengeFn(source []byte, salt []byte) ([]byte, []byte, error)
}

func New(hashtype string) (Challenge, error) {
	// split hashtype at the first hyphen
	prefix, _, _ := strings.Cut(hashtype, "-")
	switch prefix {
	case "X2R":
		return x2r.NewX2RHashType(hashtype)
	default:
		return nil, fmt.Errorf("%s is not a supported challenge type", hashtype)
	}
}
