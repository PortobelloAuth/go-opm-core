package hashing

import (
	"fmt"
	"strings"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/argon2id"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/blake3"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/noop"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/pbkdf2"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/sha2"
	"github.com/PortobelloAuth/go-opm-core/pkg/hashing/sha3"
)

type HashingConstructor func(specifier string) (Hash, error)

var hashingRegistery map[string]HashingConstructor
var newCalled = false

func init() {
	hashingRegistery = make(map[string]HashingConstructor)

	// register the built-in Hash types
	builtins := []struct {
		K string
		N func(hashtype string) (Hash, error)
	}{
		{"Argon2id", func(hashtype string) (Hash, error) { return argon2id.NewArgon2idHashType(hashtype) }},
		{"Blake3", func(hashtype string) (Hash, error) { return blake3.NewBlake3(32, nil) }},
		{"Noop", func(hashtype string) (Hash, error) { return noop.New() }},
		{"PBKDF2", func(hashtype string) (Hash, error) { return pbkdf2.NewPBKDF2HashType(hashtype) }},
		{"SHA2", func(hashtype string) (Hash, error) { return sha2.NewSHA2HashType(hashtype) }},
		{"SHA3", func(hashtype string) (Hash, error) { return sha3.NewSHA3HashType(hashtype) }},
	}

	for _, b := range builtins {
		hashingRegistery[b.K] = b.N
	}
}

func Register(hashtype string, fn HashingConstructor) error {
	// FIXME?: don't allow this to be called after New()? It makes testing annoying.
	if newCalled {
		return fmt.Errorf("Register() called after New()")
	}

	// FIXME?: don't allow overwriting builtins?
	hashingRegistery[hashtype] = fn

	return nil
}

type HashFn func(source []byte, salt []byte) ([]byte, error)
type HashFunc func(source string, salt []byte) ([]byte, error)

type Hash interface {
	HashFn(source []byte, salt []byte) ([]byte, error)
	HashFunc(source string, salt []byte) ([]byte, error)
}

func New(hashtype string) (Hash, error) {
	// TODO: make this hash lookup case insensitive
	newCalled = true

	// split hashtype at the first hyphen
	prefix, _, found := strings.Cut(hashtype, "-")
	if found {
		fn, ok := hashingRegistery[prefix]
		if ok {
			return fn(hashtype)
		}
	}

	fn, ok := hashingRegistery[hashtype]
	if ok {
		return fn(hashtype)
	}

	return nil, fmt.Errorf("%s is not a supported hash type", hashtype)
}
