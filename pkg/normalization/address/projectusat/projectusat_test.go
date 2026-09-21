package projectusat_test

import (
	"fmt"
	"testing"

	"github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/projectusat"
	"github.com/PortobelloAuth/go-projectusat/pkg/address"
	"github.com/PortobelloAuth/go-projectusat/pkg/address/parser"
)

func TestUSAtNormalizationApply(t *testing.T) {
	stub := parser.ParsingFn(func(source string) (*address.Address, error) {
		return &address.Address{
			PrimaryNumber:       "123",
			StreetName:          "Main",
			StreetSuffix:        "Street",
			SecondaryDesignator: "Apartment",
			SecondaryNumber:     "4",
			City:                "Springfield",
			Region:              "Illinois",
			Postal:              "62701",
		}, nil
	})
	expected := "123 MAIN ST # 4 SPRINGFIELD IL 62701"

	n, err := projectusat.New(stub, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := n.Apply("anything")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != expected {
		t.Errorf("output '%s' did not match expected '%s'", output, expected)
	}
}

func TestUSAtNormalizationApplyUnreadable(t *testing.T) {
	stub := parser.ParsingFn(func(source string) (*address.Address, error) {
		return nil, fmt.Errorf("no reading")
	})

	n, err := projectusat.New(stub, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := n.Apply("123 Nowhere")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != "123 Nowhere" {
		t.Errorf("output '%s' did not match expected '123 Nowhere'", output)
	}
}

func TestUSAtNormalizationApplyRejected(t *testing.T) {
	stub := parser.ParsingFn(func(source string) (*address.Address, error) {
		return &address.Address{
			PrimaryNumber:       "123",
			StreetName:          "Main",
			StreetSuffix:        "Street",
			SecondaryDesignator: "Apartment",
			SecondaryNumber:     "4",
			City:                "Springfield",
			Region:              "Illinois",
			Postal:              "62701",
		}, nil
	})
	verifier := func(a *address.Address) (*address.Address, error) {
		return nil, fmt.Errorf("rejected")
	}

	n, err := projectusat.New(stub, verifier)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	output, err := n.Apply("x")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if output != "x" {
		t.Errorf("output '%s' did not match expected 'x'", output)
	}
}

func TestUSAtNormalizationRequiresParser(t *testing.T) {
	_, err := projectusat.New(nil, nil)
	if err == nil {
		t.Errorf("expected an error when no parser is provided")
	}
}
