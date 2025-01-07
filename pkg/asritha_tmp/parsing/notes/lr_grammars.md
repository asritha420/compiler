# LR Grammars


- **LR(1) grammars**: grammars that can be parsed via shift-reduce w/ single token lookahead 
  - LR
      - L: left-to-right scanning on input
      - R: produce rightmost derivation 
  - can accommodate left recursion
  - grammar must be non-ambiguous 
  - cannot have shift-reduce or reduce-reduce conflicts


- Shift-Reduce Parsing: 
  - "bottom up" parsing strategy: start with stream of tokens, look for rules that can be applied to reduce sentential forms into non-terminals
    - for parsing to be successful, the sequence of reductions must lead to start symbol
  - *shift action*: consume one token from input stream, push it onto the stack
  - *reduce action*: apply one from of the form $A \rightarrow a$ from grammar, so replaces sentential form $a$ on stack with the corresponding non-terminal


- LR(0) automaton: all possible rules currently under consideration by a shift-reduce parser
  - also called canonical collection or compact finite state machine of grammar 

--- 

Given  $G_{10}$: 
- $P \rightarrow E$ 
- $E \rightarrow E + T$ 
- $E \rightarrow T$ 
- $E \rightarrow id(E)$
- $T \rightarrow id$ 


**1) Calculate FIRST and FOLLOW sets**

FIRST set:
- $P$: $id$
- $E$: $id$
- $T$: $id$

FOLLOW set: 

- $P$: $\$$
- $E$: $\$$, $+$, $)$
- $T$: $\$$
