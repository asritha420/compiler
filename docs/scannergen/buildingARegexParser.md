# Part 1.1: Building a Regex Parser 

#### Prerequisites 
- Regex  https://regexr.com/ 

### Introduction 

The scanner generator will take in a regex to classify each token type. 
For example, the regex for some of the token types we established in [Parts of a Scanner](partsOfAScanner.md) would be: 

- **identifier**: 
  - `[A-Z|a-z]+([A-Z|a-z]|[0-9])*`
  - *identifier is a sequence of letters and numbers, where the first character is not a number*
- **number**: 
  - `[0-9]+(.[0-9]+)?`
  - *number is a sequence numbers with an optional decimal point; if the decimal point exists, it has a sequence of numbers on either side of it* 
- **for**: 
  - `for`
- **left_paren**: 
  - `(`

Our`ScannerGen` struct will have a map field called `TokenTypes`, where the keys will store the names of all the token types, and the values will store each corresponding regex as a string.

See the `generateScanner()` in `./examples/monster/main.go` for the full `TokenTypes` map as per the Monster spec. 

These regex rules will be used in the upcoming steps of our parser. However, we must first verify that a chunk of a Monster source code follows a regex expression, and then construct an AST. #TODO: fix this explanation 

### Defining Regex 

The formal definition of regex is the following: 
#TODO: cite 
> Given an alphabet $\Sigma$, a **regular expression** (regex) $s$ is a string that denotes $L(s)$, or a set of strings drawn from $\Sigma$.  
> $L(s)$ has the following properties:  
> - For any $a \in \Sigma$, $a$ is a regex and $L(a) = \{ a \}$
> - $\epsilon$ is a regex and $L(\epsilon) = \{ \epsilon \}$  
> 
> Given two regexes $s$ and $t$:  
> - **Rule #1 (Alternation)**: $s|t$ is a regex such that $L(s|t) = L(s) \cup L(t)$  
> - **Rule #2 (Concatenation)**: $st$ is a regex such that $L(st) = L(s)L(t)$, or a string in $L(s)$ followed by a string in $L(t)$.
> - **Rule #3 (Kleene closure)**: $s*$ is a regex such that $L(s*) = L(s)$ concatenated zero or more times. 


The scanner generator will support the following syntax for the regex strings that define token types:

- **Regular Operations**: 
  - **alternation**: $a|b$
  - **concatenation**: $ab$
  - **Kleene closure**: $a*$
- **Non-regular Operations**: (can be rewritten in terms of the regular operations) 
  - **zero or one**: $a?$
    - $a$ is optional  
    - rewritten as $(a | \epsilon)$ 
  - **one or more**: $a+$
    - $a$ is repeated one or more times 
    - rewritten as $aa*$
  - **character class**: $[a-z]$
    - match any character in the range of $a$ to $z$
    - rewritten as $(a | b | ... | z)$
  - **negation character class**: $^\wedge a$
    - match any character except once 
    - rewritten as $\sum -a$

### An Unambiguous Grammar for Regex 
We will create an LL(1) grammar to interpret any regex expression. 

Lets start with grammar $G$.


Unambiguous Grammar $G_3$: 
- $P \rightarrow E$ 
- $E \rightarrow TE'$
- $E' \rightarrow | TE'$
- $E' \rightarrow \epsilon$ 
- $T \rightarrow FT'$
- 

```azure

map[byte][]string{
'P': {"E"},
'E': {"TX"},
'X': {"|TX", parsergen.Epsilon},
'T': {"FY"},
'Y': {"FY", parsergen.Epsilon},
'F': {"GM"},
'M': {"*M", parsergen.Epsilon},
'G': {"(E)", parsergen.ValidChar, parsergen.ValidInt},
		},
		[]byte{'|', '*', '(', ')'},
		[]byte{'P', 'E', 'X', 'T', 'Y', 'F', 'M', 'G'},
```
### Creating a Recursive Descent Parser for Regex 

### Adding Type Safety to our Regex Parser

#### The Expression Problem 
