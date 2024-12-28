# Grammar Rules 

### EBNF 

**EBNF (Extended Backus-Naur Form)** is the most commmon metalanguage (a language that describes other languages).

We support a modified version of EBNF with the following conventions.

<u> Naming conventions: </u> 

* non-terminals: lowercase or camelCase 
* terminals: written in double quotes (`"..."`)

<u> Production conventions: </u>

* sequence(`x y`): two elements after another with a space in between  
* alternation (`|`)
* ranges (`'[...]'`): 
  * we only support the three built in ranges: 
    * lowercase: (`[a-z]`)
    * uppercase: (`[A-Z]`)
    * numbers: (`[0-9]`)

Support for all conventions with EBNF will be incremently added. 

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
    []rune(`expressionPrime = "|" term expressionPrime | EPSILON`),
    []rune(`term = factor termPrime`),
    []rune(`termPrime = factor termPrime | EPSILON`), 
    []rune(`factor = group factorPrime`), 
    []rune(`factorPrime = "*" factorPrime | EPSILON`), 
    []rune(`group = "(" expression ")" | [a-z] | [A-Z] | [0-9]`) 
}
```




---
### Production Rules
#### Grammar
```
P -> Rule("\n"Rule)*
Rule -> WS* NT WS* "->" WS*Prod("|"Prod)*WS*
Prod -> Token*
Token -> WS*(NonTerm | Term | Range)WS*
NonTerm -> [^ \n\t\"[]]*
Term -> "\""[^]*"\""

Range -> "[" ("{"TRange"}"ICTRange | ICTRange) "]"
ICTRange -> CTRange | "^"CTRange
CTRange -> CRange | TRange
CRange -> V*
TRange -> V"-"V
```
Note: V is all characters and WS is all white space characters.

#### Ranges
Note: In order to use `^`, `-`, `{`, or `}` in a range, please escape it.

In order to easily use a range of characters in a production, you can use square brackets with a range inside like so:

`[a-z]`

This however gets clunky if you have only a lot of disjoint characters that you would like to group. For that, you can simply put a string within `[]` which will match as long as any of the characters in the string matched (similar to many alternations). Here is an example that would match "hello world":

`[helowrd]`

You can also specify to take the opposite of a range or group of characters by adding a `^` to the front. For example:

`[^A-Z] or [^helowrd]`

By default, taking the opposite will use the range $0-2^{31}$ as the language but this is often not desirable (for example with ascii encoding) and can be changed by using `{}` at the front with another range within. For example the following will match any lower case character other than g:

`[{a-z}^g]`

Note: While specifying a range for the language is possible for non inverting (not using `^`) ranges, it will not do anything and should be avoided.

Grammar for Ranges:
```
P -> "["E"]"
E -> "{"ToRange"}"InvertibleRange | InvertibleRange
InvertibleRange -> Range | "^"Range
Range -> CharRange | ToRange
CharRange -> v*
ToRange -> v"-"v
v -> all characters
```


