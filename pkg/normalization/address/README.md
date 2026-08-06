# Address Normalization

Address normalization is complicated. The appropriate standard to apply can vary
significantly based on the locality in which the address exists and the use case
for which the address is being normalized. Therefore this module provides
multiple implementation options.

## Options

- `projectusat` - (not implemented yet) employs
  [go-projectusat](https://github.com/PortobelloAuth/go-projectusat) which
  implements the
  [Project US@ standard](https://asapnet.org/wp-content/uploads/2022/03/Project_US_FINAL_Technical_Specification_Version_1.0.pdf)
  (which is an extension of the USPS publication 28 standard) directly in Go
- `libpostal` (partially implemented) an http API wrapper around
  [pelias/libpostal-service](https://hub.docker.com/r/pelias/libpostal-service/tags)
  to normalize international addresses using libpostal, which is derrived from
  Open Street Map data

You are required to Register the address normalization implementation that you
want to use; there is no default at this time.

```go
import (
  "github.com/PortobelloAuth/go-opm-core/pkg/normalization/address/libpostal"
  "github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

func init() {
  lpEndpoint, found := os.LookupEnv("LIBPOSTAL_ENDPOINT")
  // normalization.Register() must be called before normalization.New()
  normalization.Register(
    "address",
    func(specifier string) (Normalization, error) {
      return libpostal.NewService(
        lpEndpoint,
        300 * time.Millisecond
      )
    },
    false,
  )
}
```

## Notes

Open Person Matching generally wants to treat an address as a complete formatted
string, similar to what may appear in the `"formatted"` field of an OpenID
Connect ID Token's `"address"` field. This aligns well with the
`normalization.Normalization` interface's `Apply(string) (string, error)`
method. However, there are likely serveral use cases where it is advantageous to
work with parsed address data. It has not yet been determined how we want to
handle these cases.
