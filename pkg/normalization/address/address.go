package address

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var implementation func(string) (AddressNormalization, error)
var timedout = false
var registered = make(chan bool, 1)

func init() {
	err := CheckMode()
	if err != nil {
		panic(err)
	}
}

func CheckMode() error {
	// check the
	mode, found := os.LookupEnv("GO_OPM_CORE_ADDRESS_NORMALIZATION_MODE")
	if !found {
		mode = "interface"
	}

	switch mode {
	case "projectusat":
		// implementation = projectusat.New
		return fmt.Errorf("%s address normalization mode is not yet implemented", mode)
	case "libpostal":
		// implementation = gopostal.New
		return fmt.Errorf("%s address normalization mode is not yet implemented", mode)
	case "interface":
		// do nothing. New() will error until RegisterAddressNormalizationInterface has been called.
		return nil
	default:
		return fmt.Errorf("%s is not a recognized GO_OPM_CORE_ADDRESS_NORMALIZATION_MODE", mode)
	}
}

func DeregisterAddressNormalizationInterfaceForTests() {
	// TODO: consider whether this exported test helper is how we really want to do this
	if testing.Testing() {
		implementation = nil
		close(registered) // close the previous channel in case it wasn't alredy closed
		registered = make(chan bool, 1)
		timedout = false
		return
	}
	panic("DeregisterAddressNormalizationInterfaceForTests() must only be called from tests")
}

func RegisterAddressNormalizationInterface(impl func(string) (AddressNormalization, error)) error {
	if implementation != nil {
		return fmt.Errorf("Address normalization implementation already registered")
	}
	implementation = impl
	close(registered)
	return nil
}

type AddressNormalization interface {
	Apply(input string) (string, error)
}

func New(specifier string) (AddressNormalization, error) {
	// wait a short time (once) for registration to occur
	if implementation == nil && !timedout {
		select {
		case <-registered:
		case <-time.After(1 * time.Second):
			timedout = true
		}
	}
	if implementation == nil {
		return nil, fmt.Errorf("no registered implementation; did you call RegisterAddressNormalizationInterface()?")
	}
	return implementation(specifier)
}
