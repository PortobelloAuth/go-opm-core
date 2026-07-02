package hashing

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/lib/hashing/argon2id"
	"github.com/PortobelloAuth/go-opm-core/lib/hashing/blake3"
	"github.com/PortobelloAuth/go-opm-core/lib/hashing/noop"
	"github.com/PortobelloAuth/go-opm-core/lib/hashing/pbkdf2"
	"github.com/PortobelloAuth/go-opm-core/lib/hashing/sha2"
	"github.com/PortobelloAuth/go-opm-core/lib/hashing/sha3"
)

type HashFn func(source []byte, salt []byte) ([]byte, error)
type HashFunc func(source string, salt []byte) ([]byte, error)

type Hash interface {
	HashFn(source []byte, salt []byte) ([]byte, error)
	HashFunc(source string, salt []byte) ([]byte, error)
}

func New(hashtype string) (Hash, error) {
	// split hashtype at the first hyphen
	prefix, _, _ := strings.Cut(hashtype, "-")
	switch prefix {
	case "Argon2id":
		return argon2id.NewArgon2idHashType(hashtype)
	case "Blake3":
		return blake3.NewBlake3(32, nil)
	case "Noop":
		return noop.New()
	case "PBKDF2":
		return pbkdf2.NewPBKDF2HashType(hashtype)
	case "SHA2":
		return sha2.NewSHA2HashType(hashtype)
	case "SHA3":
		return sha3.NewSHA3HashType(hashtype)
	default:
		return nil, fmt.Errorf("%s is not a supported hash type", hashtype)
	}
}
