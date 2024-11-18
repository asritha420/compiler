# TODO: 

## scanner-gen 

The scanner-generator package will read in a specification file with regex rules and produce a scanner. 

**Tasks**:
1) Input specification: 
An input file will contain a list of rules in the following format: `NAME REGEX_PATTERN`

Example: 
```ispec
NUMBER [0-9]+
IDENT     [a-zA-Z_][a-zA-Z0-9_]*
STRING    ".*?"
```

2) Convert input specifications into a parse tree
3) Write regex parser manually
5) Convert syntax tree into NFA 
6) Convert NFA to DFA 
7) Optimize DFA? 
7) Generate scanner from DFA
   8) Generate error handler for unmatched input/invalid tokens
9) Unit tests for each phase 
10) Printer
11) CLI to allow users to run the following command: 
Ex. `go run scanner-gen.go --spec lexer.spec --output lexer.go`
12) Finish `scanner-gen.md`
13) include support for more than digits and lower case letters, what about capital letters??? also include suppport for more regex rules

TODO: before starting the blog post, write out the rest of the stuff I need to add to the grammar to add support for all regex rules

14) is there a better way to generate the first and follow sets thats more efficient 

## parser-gen
1) Error handling in Crafting Interpreters parser, finish implementation  
2) Create Parser generator 
2) Finish Follow function implementation 
3) Implement memoization for first and follow
3) Port over regex parser from scanner to be generated from parser-gen
5) object algebras 
6) Include support for X' as rule name in a grammar object?
   7) Include rules as an object to allow support for capital letters in the regex grammar 
7) When initializing a grammar object, include checks to see if its ambiguous, then throw error if so 
8) more efficient way to implement generateFirstSet() and generateFollowSet()?
9) current implementation of type Expression is not type safe, also for regexExpr as well
10) in regexParser.go, I have an AST printer, does that make sense? is there a simpler way using prebuild to shit to print the nodes hierarchically. either way, implement the same thing for the main parser generator as a debugging tool
 
## general

- Add some sort of documentation validator for the functions
- Add tests for each step 
  - recursive parser generator 
  - grammar first and follow 
  - scannar helper methods 
  - etc

## Blog outline: 
- Introduction //describe parts of compiler, intention behind project 
- Building a Scanner Generator 
  - Building a Regex Parser Part 1 
  - Building a Regex Parser Part 2 //include support for the regex operations, can this be combined into one part. how long should each blog post be? 
  - Adding Type Safety to the Regex Parser
  - Coding a Finite State Machine 
  - Putting it together //Adding input building blocks for the scanner generator ie. token class, etc. need to split this up  
  - Adding Error Handling


When writing out the blog post, add functions that need to be scaffolded out at a later date, ie. checking for if a valid grammar