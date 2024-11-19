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


We will support the following regex syntax

### Defining a Grammar for Regex 
Lets start with grammar $G$. 






