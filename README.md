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

This project is licensed under the Apache License 2.0 - see the
[LICENSE](LICENSE.md) file for details.

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

### Network Settings

An Open Person Matching network is defined by its
[Network Settings](./pkg/network_settings/). These settings define how
requestors, responders, and facilitators (aka Nodes) interact to securely match
and exchange demographic hashes and credentials. Because nodes may participate
in multiple networks at once, Network Settings enable the gradual adoption, use,
and deprecation of network settings such as versions, network types, matching
rules, node/network authentication, encryption, etc. This is a critical
capability to maintain the security of an open person matching network.

### Rules

[Rules](./pkg/rules/) define the matching criteria used by a network. They are
part of a network's settings. At a high level, network rules have a root RuleSet
that contains Rules; these Rules may contain RuleDemographics and other
RuleSets. This creates a tree-like data structure that allows a complex set of
rules to be defined relatively simply. RuleDemographics and RuleSets each may
define Normalizations to be applied to their respective demographic values.
RuleDemographics may also define DemographicAlternatives functions to enable the
use of fuzzy matching, nicknames, etc.

### Normalizations

[Normalizations](./pkg/normalization/) transform demographics in defined ways to
establish a standardized form of the demographic that is suitable for
comparison. Some Normalizations are relatively simple - they convert to
uppercase, remove whitespace, etc. Others, like
[`address`](./pkg/normalization/address/) are quite complicated; by default
[`address`](./pkg/normalization/address/) uses
[`projectusat`](./pkg/normalization/address/projectusat/) which implements the
[Project US@ standard](https://asapnet.org/wp-content/uploads/2022/03/Project_US_FINAL_Technical_Specification_Version_1.0.pdf)
standard using
[`go-projectusat`](https://github.com/PortobelloAuth/go-projectusat), a Go
module that we developed for this purpose.

### Demographic Alternatives

Many demographics can show up in a variety of forms in a person's records:
diacritics in a person's name or address, nick names, misspellings, historical
names, etc. In some cases these are represented by distinct entries in a raw
demographic map. For others a
[Demographic Alternatives Function](./pkg/demographicalts/) provides a mechanism
for Rules to define different ways that a particular demographic may be
considered to match.

## TODO

- implement a RuleSet corresponding to the current
  [CMS HealthTech Ecosystem
  patient matching rules](https://patientmatching.org/)
- implement the hash-routed-challenge network type and define the (probable)
  `Request()`, `Register()`, and `Respond()` APIs
- implement other meaningful rule sets referenced in TEFCA SoPs and HL7 FHIR
  $match standards
- add a network authentication type package
- add asymmetric encryption type package
- add symmetric encryption package
- revisit the `match.Demographics` type; is it useful as it is?
