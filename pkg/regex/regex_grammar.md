TODO: move this file
// TODO: what should be public and what should not in the parser
// TODO: custom error
// TODO: add escape characters 
```
Regex         -> Alt

Alt           -> Concat AltPrime
AltPrime      -> "|" Concat AltPrime | ε

Concat        -> Repeat ConcatPrime
ConcatPrime   -> Repeat ConcatPrime | ε

Repeat        -> Group Quantifier?
Quantifier    -> "*" | "+" | "?"

Group         -> "(" Regex ")" | CharRange | Char

CharRange     -> "[" CharRangeBody "]"
CharRangeBody -> "^"? (CharRangeAtom)+   
CharRangeAtom -> Char ("-" Char)?       

Char          -> ANY_VALID_CHAR
```

```
note that v represents any valid char 
FIRST SETS: 
Regex: (, [, v
Alt: (, [, v
AltPrime: | 
Concat: (, [, v
ConcatPrime: (, [, v
Repeat: (, [, v
Quantifier: *, +, ? 
Group: (, [, v
CharRange: [
CharRangeBody: ^, v
CharRangeAtom: v 
Char: v
```
