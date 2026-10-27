# Rules and RuleSets

Open Person Matching networks employ rule sets that define the matching criteria
for the network. A RuleSet is a tree-like graph data structure with a root
RuleSet containing Rules and, optionally, Normalizations to be applied to those
Rules. A Rule contains Demographics that it acts on. A Rule may also contain
another RuleSet; this enables the composition of Rules.

- Network settings adopt RuleSets
- RuleSets are made up of matching Rules
- Rules may contain a list of RuleDemographics or a RuleSet, not both
- If any Rule in a RuleSet matches it is considered a match
- All RuleDemographics must match for a Rule to match
- RuleSets may employ global Normalizations that are applied to all
  RuleDemographics before their own Normalizations apply
- Rules must encode:
  - Normalizations employed
  - Fuzzy matching logic allowed
  - Source demographic data
- A Rule compares the output of the demographic data, with normalizations
  applied, split in to a fuzzy matching set, and then used to generate a hash or
  challenge
