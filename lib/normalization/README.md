# Normalization

This package includes several normalization functions. These functions take a string
input and transform it in a defined manner in order to reduce variations between
demographic values that should otherwise match.

Examples include:

- `lowercase` - all letters in the identifier should be converted to lowercase
- `uppercase` - all letters in the identifier should be converted to uppercase
- `s/{1}/{2}` - all sequences matching the pattern {1} should be replaced with
  the sequence {2}. For this version of this specification patterns should be
  syntactically correct Perl-compatible Regular Expressions though only very
  simple substitutions are supported (i.e. look-ahead, look-behind, capture
  groups, anchors, etc. will result in the Facilitator rejecting the Request)
- `no_space` - whitespace is removed
- `no_punct` - punctuation, such as `.`, `,`, `#`, `-`, etc. is removed
- `us-101-key` - the conversion of all letters to suitable lowercase English
  letters through lowercasing and typographic approximation (needs explicit
  specification) to letters which may be input using a standard US 101-key
  keyboard without the use of any modifier keys
- `us-reduced_v0` - a standardized combination of normalization rules intended
  to be used as a default for names and addresses used in the United States of
  America consisting of: `no_punct`, `no_space`, and `us-101-key`. While it will
  not cover all cases, it is intended to account for many of the variations that
  are likely to appear in demographic data that has been digitally recorded by
  humans using a keyboard.