# Grammar Rules 

### EBNF 

**EBNF (Extended Backus-Naur Form)** is the most commmon metalanguage (a language that describes other languages).

We support a modified version of EBNF with the following conventions.

<u> Naming conventions: </u> 

* non-terminals: lowercase or camelCase 
* terminals: written in double quotes (`"..."`)

<u> Production conventions: </u>
* the last production is used as the start production?

| Function | Notation |
|----------|----------|
| definition | = |
| concatenation | space |
| termination | ; |
| alternation | \| |
| exception | - |
| none or one | ? |
| none or more | * |
| one or more | + |
| grouping | (...) |
| terminal | "..." |
| range | [...] |
| escape | & |

Note: Any number of space characters (`\t, \n, \v, \f, \r, ' ', 0x85, 0xA0`) will be accepted for concatenation or between any functions but there should either be at most 1 space.

Support for all conventions with EBNF will be inclemently added. 


### Defining grammar 

First, the list of valid terminals and non-terminals must be determined.

Some common terminals include: 
* identifiers: names for variables, classes, functions, etc. 
* keywords
* literals: string literals, numeric literals, boolean literals, etc. 
* separators and delimiters: {, }, (, ), ;, etc.

Any grammar is a list of production rules. Each production rule will contain a non-terminal, and the "formula" for composing that non-terminal in a valid structure.  

There are two options two define the grammar of a language and initialize the parser generator. 

The examples below show how to define a regex grammar in the two supported formats: 

**Option 1: Text File Based**

```


```

**Option 2: Programmatically create an instance of `parser.grammar`**

```
grammar := [][]rune{ 
    []rune(`production = expression`), 
    []rune(`expression = term expressionPrime`),
    []rune(`expressionPrime = "|" term expressionPrime | ""`),
    []rune(`term = factor termPrime`),
    []rune(`termPrime = factor termPrime | ""`), 
    []rune(`factor = group factorPrime`), 
    []rune(`factorPrime = "*" factorPrime | ""`), 
    []rune(`group = "(" expression ")" | [a-z] | [A-Z] | [0-9]`) 
}
```


---
### Production Rules
#### Grammar (Using Go String formatting)
```
rangeChar = ([] - [&-&&&[&]]) | "&&-" | "&&&&" | "&&[" | "&&]";identifierChar = [a-z] | [A-Z] | [0-9] | "_";
stringChar = ([] - ["&&]) | "&&&"" | "&&&&";
spaceChar = [\t\n\v\f\rU+0085U+00A0];

terminal = "&"" stringChar* "&"";
identifier = identifierChar+;
toRange = (rangeChar "-" rangeChar);
charRange = rangeChar*;
range = "[" toRange | charRange "]";

space = spaceChar*;

term = "(" RHS ")"
     | terminal
     | identifier
     | range;

factor = term (space ([?*+] | ("-" space term)))?;

concatenation = factor (space factor)*;
alternation = concatenation (space "|" space concatenation)*;

RHS = space alternation space;
LHS = space identifier space;

rule = LHS "=" RHS ";";

grammar = rule*;
```

#### Ranges
Note: In order to use `-` or `&` in a range, please escape it.

In order to easily use a range of characters in a production, you can use square brackets with a range inside like so:

`[a-z]`

This however gets clunky if you have only a lot of disjoint characters that you would like to group. For that, you can simply put a string within `[]` which will match as long as any of the characters in the string matched (similar to many alternations). Here is an example that would match "hello world":

`[helowrd ]`

There is also a short hand of an empty range `[]` which defines all characters `0-2^32`
