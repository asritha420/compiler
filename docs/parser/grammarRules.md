# Grammar Rules 

### EBNF 

**EBNF (Extended Backus-Naur Form)** is the most commmon metalanguage (a language that describes other languages).

We support a modified version of EBNF with the following conventions.

<u> Naming conventions: </u> 

* non-terminals: lowercase or camelCase 
* terminals: written in double quotes (`"..."`)

<u> Production conventions: </u>

* sequence: two elements after another 
* alternation (`|`): #TODO 
* grouping: (`(...)`): #TODO 
* 

Support for all conventions with EBNF will be incremently added. 

### Defining grammar 

First, the list of valid terminals and non-terminals must be determined.

Some common terminals include: 
* identifiers: names for variables, classes, functions, etc. 
* keywords
* literals: string literals, numeric literals, boolean literals, etc. 
* separators and delimiters: {, }, (, ), ;, etc.

Terminals must be enclosed in "". 

Any grammar is a list of production rules. Each production rule will contain a non-terminal, and the "formula" for composing that non-terminal in a valid structure.  

There are two options two define the grammar of a language and initialize the parser generator. 

The examples below define the same grammar in the two different formats: 

**Option 1: Text File Based**

```


```


**Option 2: Programmatically create an instance of `parser.grammar`**




