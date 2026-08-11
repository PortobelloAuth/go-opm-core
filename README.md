# Go Open Person Matching Core

Open Person Matching is designed as a privacy preserving person matching
algorithm employing:

1. the exchange of hashes of demographic data to determine whether Requestors
   and Responders have data about the same person
2. network settings that define how demographics are exchanged, what matching
   rules are employed, how data is normalized, etc.
3. network settings also allow for inherent upgradability as Requestors,
   Responders, and Facilitators participate in multiple networks

This library is focused on the core patient matching interfaces related to Open
Person Matching, without the operational concerns of Requestors, Responders, and
Facilitators. Each package is intended to isolate and address an important,
focused concern.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE.md) file for details.

## Concepts

The top level package exposes two functions, `Match()` and `PlaintextMatch()`.
`Match()` is this library's primary entrypoint. It compares a map of plain text
string identifiers against a set of hashed identifiers that claim to have been
generated using a given hashing algorithm, salt, encoding function, and matching
rules (which incorporate normalization functions.) `Match()` applies the hashing
algorithm, salt, encoding function, and matching rules to the plain text
identifiers and checks to see if any of the resulting hashed identifiers match
the supplied hashed identifiers. If a single matching rule matches, the result
will be true.

## TODO

- implement a working address parser and normalizer (using go-projectusat)
- implement a RuleSet corresponding to the current CMS HealthTech Ecosystem
  patient matching rules
- implement the hash-routed-challenge network type and define the (probable)
  `Request()`, `Register()`, and `Respond()` APIs
- implement other meaningful rule sets referenced in TEFCA SoPs and HL7 FHIR
  $match standards
- add a network authentication type package
- add asymmetric encryption type package
- add symmetric encryption package
- revisit the `match.Demographics` type; is it useful as it is?
- finish / fix address normalization
