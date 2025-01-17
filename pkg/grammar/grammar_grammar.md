Tokens for Scanner
```
letter = [a-z] | [A-Z]
digit = [0-9]
space = \s+
```

Grammar Rules
```
start = rules;

symbol = "[" | "]" | "{" | "}" | "(" | ")" | "<" | ">"
       | "'" | "\"" | "=" | "|" | "." | "," | ";" | "-" 
       | "+" | "*" | "?" | "_" | "/" | "\\";

strChar = "letter" | "digit" | "space" | symbol ;
idChar = "letter" | "digit" | "_";

identifier = idChar | idChar, identifier;
string = epsilon | strChar, string;
token = "\"" string "\"";
separator = space | epsilon;

term = token | identifier | "(", rhs, ")";
sTerm = separator, term, separator;

unary = "?" | "*" | "+" | epsilon
factor = sTerm, unary, separator;

concatenation = factor | factor, ",", concatenation;
alternation = concatenation | concatenation, "|", alternation;

lhs = identifier;
rhs = alternation;

rule = separator, lhs, separator, "=", rhs, ";", separator;
rules = rule | rule, rules;
```