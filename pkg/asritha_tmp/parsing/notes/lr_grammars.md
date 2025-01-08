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


- *LR(0) automaton*: represents choices available at any step of bottom up parsing  
  - also called canonical collection or compact finite state machine of grammar
  - algorithm:
    - Create State 0
      - Get Kernel of State 0 
      - Get Closure of State 0: for each item in the state with a non-terminal X immediately to the right of the dot, add all rules in the grammar that have X as a non-terminal
      - Create transitions for each of the symbols to the right of the dot 
    - For each transition, create a new state containing the matching items with dot moved one position to the right
    - Continue procedure until no new items can be added
    - Key points: 
      - kernel items: start items + items where dot is not at the beginning 
      - non-kernel items: items introduced by the closure operation 
  - state w/ item w/ dot at end rule: possible reduction 
  - transition on a terminal that moves dot one position to the right: possible shift
  - conflicts: 
    - *shift-reduce conflict*: choice between shift action and reduce action in the same state 
    - *reduce-reduce conflict*: two distinct rules have been matched 

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
