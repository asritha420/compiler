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
