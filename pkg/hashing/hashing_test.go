package hashing_test

import (
	"encoding/hex"
	"slices"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/hashing"
)

func TestBlake3HashFn(t *testing.T) {
	expected, err := hex.DecodeString("c6d25f73c520a48e5c193586a2265574d2278be2c16f6a8d91c8f087e4ba1b21")
	if err != nil {
		t.Errorf("Error decoding expected value for testing New(\"Blake3\"): %s", err)
	}

	h, err := hashing.New("Blake3")
	if err != nil {
		t.Errorf("Error creating New(\"Blake3\"): %s", err)
	}

	input := []byte("input string")
	salt := []byte("salt for testing")
	output, err := h.HashFn(input, salt)
	if err != nil {
		t.Fatalf("Error generating hash with blake3.HashFn(): %s", err)
	}

	if !slices.Equal(expected, output) {
		t.Errorf("HashFn output byte slice '%x' does not match expected slice %x", output, expected)
	}

	output2, err := h.HashFunc(string(input), salt)
	if err != nil {
		t.Errorf("Error generating hash with blake3.HashFunc(): %s", err)
	}

	if !slices.Equal(expected, output2) {
		t.Errorf("HashFunc output byte slice '%x' does not match expected slice %x", output2, expected)
	}
}

func TestNoopHashFn(t *testing.T) {
	h, err := hashing.New("Noop")
	if err != nil {
		t.Errorf("Error creating hashing.New(\"Noop\"): %s", err)
	}

	input := []byte("input string")
	salt := []byte("test salt")
	output, err := h.HashFn(input, salt)
	if err != nil {
		t.Fatalf("Error running noop.HashFn(): %s", err)
	}

	if !slices.Equal(input, output) {
		t.Errorf("Input byte slice '%s' does not match output slice '%s'", input, output)
	}
}

var builtinTestCases = []struct {
	Spec string
	Hex  string
}{
	{"Argon2id-1-20mb-256", "d3c246d31c57369c99116e6c14ba9db041ad532e2595d2a1ebe46ce0b15134c8"},
	{"Argon2id-1-64mb-256", "fab644f354ac42c6ce9c65b4952a6f50818661c9526cef5ee6b860b37e5f5b00"},
	{"Blake3", "fb1678675d1a49a589137ab17c35dc3150a8ab445d4190bc78346b77da91a380"},
	{"Noop", "74657374"}, // "test" in hex encoding
	{"PBKDF2-SHA2_256-100000-256", "83726f70443a10c00d12141a1f6a96f519606f4e123f4c7f18cb19d31ed4aa91"},
	{"SHA2-224", "0b886dcc1117156500446d323698ac9e2a883e2846e1c56448375ec7"},
	{"SHA2-256", "b197d2241c1e9ccecf285e8ab4a5b578bf7c2c07ea2fc0285b8efceeb42ba63c"},
	{"SHA2-384", "84485bec48fcca91edfd803fc08c4e4b1c7353e1f8680ea4cafcbc643998f3da37773d01af0e9e49e140650cad8f7c23"},
	{"SHA2-512", "3775b3f6f0db255957fa6e8c1f1fcff38bc62a0df1ec76ab043a3811dd1af98467bd7b12ff3614d89ba43d9e63bba9cf59b52a715b66f8c8a6462da65653ae7a"},
	{"SHA2-512_224", "e0c707d11c9fa1b038329c847521a6b699e6f76e425d81242eb35358"},
	{"SHA2-512_256", "e8145a2d5ddca2eb263f3ee843b98f5f485014d98f57013b1d06f39005494276"},
	{"SHA3-224", "b836c5675735ffd031d2c45a8f22793dfab8a0539f4c8073f060c00f"},
	{"SHA3-256", "551250bb5fa2e030a603eef2166a93ac35ed3d88d841e62a8ee53edbb3e5169b"},
	{"SHA3-384", "b7fb0f384f037a6c08cf2bb2dc1285ed1f872c32b3c2ce15794b5416923755e94ec1f507200ede51a93fe1cbbbd5d65e"},
	{"SHA3-512", "08f72e376ed90daac4fd49a5d6683f2f4109abb9d0331240a865bdfa9cade40beba1119c97a44fe932de55fec98d4331d29c70b88c13aad8dca41e790b7a5e32"},
}

func TestBuiltInHashFns(t *testing.T) {
	originalSalt := []byte("salt for testing")
	for _, tc := range builtinTestCases {
		t.Run(tc.Spec, func(t *testing.T) {
			expected, err := hex.DecodeString(tc.Hex)
			if err != nil {
				t.Errorf("Error preparing to test New(%q): %s", tc.Spec, err)
			}

			h, err := hashing.New(tc.Spec)
			if err != nil {
				t.Fatalf("Error creating New(%q): %s", tc.Spec, err)
			}

			input := []byte("test")
			salt := slices.Clone(originalSalt)
			output, err := h.HashFn(input, salt)
			if err != nil {
				t.Errorf("Error generating hash with h.HashFn(): %s", err)
			}
			if string(salt) != string(originalSalt) {
				t.Errorf("%q h.HashFn() modified provided salt", tc.Spec)
			}

			if !slices.Equal(expected, output) {
				t.Errorf("HashFn output byte slice '%x' does not match expected slice '%x'", output, expected)
			}

			output2, err := h.HashFunc(string(input), salt)
			if err != nil {
				t.Errorf("Error generating hash with blake3.HashFunc(): %s", err)
			}
			if string(salt) != string(originalSalt) {
				t.Errorf("%q h.HashFn() modified provided salt", tc.Spec)
			}

			if !slices.Equal(expected, output2) {
				t.Errorf("HashFunc output byte slice '%x' does not match expected slice '%x'", output2, expected)
			}
		})
	}
}
