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
S = space | epsilon;

term = token | identifier | "(", rhs, ")";
sTerm = S, term, S;
concatenation = sTerm | sTerm, ",", concatenation;
alternation = concatenation | concatenation, "|", alternation;

lhs = identifier;
rhs = alternation;

rule = S, lhs, S, "=", rhs, ";", S;
rules = rule | rule, rules;
```