# Rules and RuleSets

Open Person Matching networks employ rulesets that define how matching criteria
are

- Network settings adopt RuleSets
- RuleSets are made up of matching Rules
- Rules may contain a list of RuleDemographics or a RuleSet, not both
- If any Rule in a RuleSet matches it is considered a match
- All RuleDemogrphics must match for a Rule to match
- RuleSets may employ global Normalizations that are applied to all
  RuleDemographics before their own Normalizations apply
- Rules must encode:
  - Normalizations employed
  - Fuzzy matching logic allowed
  - Source demographic data
- A Rule compares the output of the demographic data, with normalizations
  applied, split in to a fuzzy matching set, and then used to generate a hash or
  challenge
