Tokens for Scanner
```
letter = [a-z] | [A-Z]
digit = [0-9]
space = (" " | "\n" | "\t" | "\r" | "\f" | "\b")*
```

Grammar Rules
```
start = rules

symbol = "[" | "]" | "{" | "}" | "(" | ")" | "<" | ">"
       | "'" | '\"' | "=" | "|" | "." | "," | ";" | "-" 
       | "+" | "*" | "?" | "\n" | "\t" | "\r" | "\f" | "\b" ;

character = letter | digit | symbol | "_" | " " ;

identifierChar = letter | digit | "_"
identifier = identifierChar | identifierChar, identifier
chars = epsilon | character, chars
terminal = "\"" chars "\""

term = terminal | identifier | "(", rhs, ")"
concatenation = term | term, ",", concatenation
alternation = concatenation | concatenation, "|", alternation

lhs = identifier
rhs = alternation

rule = lhs, "=", rhs, ";"
rules = rule | rule, rules
```