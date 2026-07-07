# Address Normalization

Address normalization is complicated. The appropriate standard to apply can vary significantly based on the locality in which the address exists adn the use case for which the address is being normalized. Therefore this module provides multiple implementation options.

## Options

- `projectusat` - (not implemented yet) implements the [Project US@ standard](https://asapnet.org/wp-content/uploads/2022/03/Project_US_FINAL_Technical_Specification_Version_1.0.pdf) (which is an extension of the USPS
  publication 28 standard) directly in Go
- `gopostal` (not implemented yet, requires libpostal to be installed) uses the gopostal module to normalize
  international addresses using libpostal, which is derrived from Open Street Map data
- `interface` (not implemented yet) allows registering a constructor for an implementation of the
  `address.AddressNormalization` interface. This enables the use of tools like the USPS Address normalization
  API and other address normalization APIs.

## Notes

Open Person Matching generally wants to treat an address as a complete formatted string, similar to what may
appear in the `"formatted"` field of an OpenID Connect ID Token's `"address"` field. This aligns well with
the `normalization.Normalization` interface's `Apply(string) (string, error)` method. However, there are
likely serveral use cases where it is advantageous to work with parsed address data. It has not yet been
determined how we want to handle these cases.
