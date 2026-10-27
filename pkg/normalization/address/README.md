# Address Normalization

Address normalization is complicated. Addressing rules vary significantly based
on the locality in which the address exists and the use case for which the
address is being normalized. Address normalization currently defaults to our
`projectusat` package. This employs
[go-projectusat](https://github.com/PortobelloAuth/go-projectusat), an address
parsing and normalization library that implements the
[Project US@ standard](https://asapnet.org/wp-content/uploads/2022/03/Project_US_FINAL_Technical_Specification_Version_1.0.pdf)
(which is an extension of the USPS publication 28 standard) directly in Go.
`go-projectusat` employs TIGER file data from the US Census Bureau encoded as
bloom filters to enable data dependent normalization for US addresses with a
relatively small memory footprint and no API calls.

## Options

Did we mention that address normalization is complicated? Your use case may
require significantly different normalization rules than those prescribed by the
[Project US@ standard](https://asapnet.org/wp-content/uploads/2022/03/Project_US_FINAL_Technical_Specification_Version_1.0.pdf).
Therefore this module allows the caller to override the default implementation.

To employ a different address normalization implementation you must register it
with `normalization.Register()` before `normalization.New()` is first called.

```go
import (
  ".../example"
  "github.com/PortobelloAuth/go-opm-core/pkg/normalization"
)

func init() {
  endpoint, found := os.LookupEnv("EXAMPLE_ENDPOINT")
  // normalization.Register() must be called before normalization.New()
  normalization.Register(
    "address",
    func(specifier string) (Normalization, error) {
      return example.NewService(
        endpoint,
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
