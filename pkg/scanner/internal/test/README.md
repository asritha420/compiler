# EBNF Grammar

We will start with the original EBNF grammar from wikipedia which closely follows ISO/IEC 14977

```
letter = "A" | "B" | "C" | "D" | "E" | "F" | "G"
       | "H" | "I" | "J" | "K" | "L" | "M" | "N"
       | "O" | "P" | "Q" | "R" | "S" | "T" | "U"
       | "V" | "W" | "X" | "Y" | "Z" | "a" | "b"
       | "c" | "d" | "e" | "f" | "g" | "h" | "i"
       | "j" | "k" | "l" | "m" | "n" | "o" | "p"
       | "q" | "r" | "s" | "t" | "u" | "v" | "w"
       | "x" | "y" | "z" ;

digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;

symbol = "[" | "]" | "{" | "}" | "(" | ")" | "<" | ">"
       | "'" | '"' | "=" | "|" | "." | "," | ";" | "-" 
       | "+" | "*" | "?" | "\n" | "\t" | "\r" | "\f" | "\b" ;

character = letter | digit | symbol | "_" | " " ;
identifier = letter , { letter | digit | "_" } ;

S = { " " | "\n" | "\t" | "\r" | "\f" | "\b" } ;

terminal = "'" , character - "'" , { character - "'" } , "'"
         | '"' , character - '"' , { character - '"' } , '"' ;

terminator = ";" | "." ;

term = "(" , S , rhs , S , ")"
     | "[" , S , rhs , S , "]"
     | "{" , S , rhs , S , "}"
     | terminal
     | identifier ;

factor = term , S , "?"
       | term , S , "*"
       | term , S , "+"
       | term , S , "-" , S , term
       | term , S ;

concatenation = ( S , factor , S , "," ? ) + ;
alternation = ( S , concatenation , S , "|" ? ) + ;

rhs = alternation ;
lhs = identifier ;

rule = lhs , S , "=" , S , rhs , S , terminator ;

grammar = ( S , rule , S ) * ;
```

This grammar uses the following notation for its functions

| Function | Notation |
|----------|----------|
| definition | = |
| concatenation | , |
| termination | ; or . |
| alternation | \| |
| none or one | [...] or ? |
| none or more | {...} or * |
| one or more | + |
| grouping | (...) |
| terminal | "..." or '...' |
| special sequence | ?...? |
| exception | - |

This standard was meant to be very extensible so it has the `special sequence` function but we do not need it to be as flexible so we can simply remove it. Also, we want to use the entire unicode character space so defining terminals by single characters will be impractical (just look at how long letter is). We will add ranges to solve this problem. Much like regex, we will use `[]` in which you can either define a string of characters or a range defined by 2 characters separated by `-`. This conflicts with `none or one` so we will remove `[]` from its notation. We also only need 