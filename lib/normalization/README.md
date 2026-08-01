# Normalization

This package includes several normalization functions. These functions take a
string input and transform it in a defined manner in order to reduce variations
between demographic values that should otherwise match.

Examples include:

- `lowercase` - all letters in the identifier should be converted to lowercase
- `uppercase` - all letters in the identifier should be converted to uppercase
- `s/{1}/{2}` (aka `replace`) - all sequences matching the pattern {1} should be
  replaced with the sequence {2}. This is a simplistic find and replace match
  that replaces matching sequences as it finds them.
- `whitespace` - whitespace is removed
- `punctuation` - punctuation, such as `.`, `,`, `#`, `-`, etc. is removed
- `transliterate` - replace all unicode runes with reasonable ASCII
  approximations as implemented by
  [anyascii](https://github.com/anyascii/anyascii)
- `us-101-key` - all unicode runes are converted to suitable lowercase English
  letters through the `transliterate` and `lowercase` normalizations. This
  appoximates how these runes may be input using a standard US 101-key keyboard
  without the use of any modifier keys.
- `us-reduced_v0` - a standardized combination of normalization rules intended
  to be used as a default for names and addresses used in the United States of
  America consisting of: `no_punct`, `no_space`, and `us-101-key`. While it will
  not cover all cases, it is intended to account for many of the variations that
  are likely to appear in demographic data that has been digitally recorded by
  humans in the United States using a keyboard.
