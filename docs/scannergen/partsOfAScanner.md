# Parts of a Scanner 

## Tokens
A chunk of raw Monster source code can be classified into a token.

Based on the `Monster` spec, we will have the following token types:

**Single Character**:
1) left_paren
2) right_paren
3) left_brace
4) right_brace
5) comma
6) dot
7) minus
8) plus
9) semicolon
10) slash
11) star

**Single or Double Character**:
1) bang
2) bang_equal (`!=`)
3) equal (`=`)
4) equal_equal
5) greater
6) greater_equal
7) less
8) less_equal

#TODO: double_slash for comments?  

**Literals**:
1) identifier
2) string
3) number

**Keywords**:
1) var
2) const
3) if
4) else
5) func
6) return
7) for
8) int
9) bool
10) string
11) char
12) print

For example, the `.mon` program below will produce the following stream of tokens: `func`, `identifier`, `left_paren`, `int`, `identifier`, `bool`, `left_brace`, `for`, `left_paren`, `int`, `identifier`, `equal`, `number`, `semicolon`, `identifier`, `less`, `semicolon`, `identifer`, `plus`, `plus`, `right_paren`, `left_brace`, `print`, `left_paren`, `identifier`, `right_paren`, `right_brace`, `return`, `left_paren`, `identifier`, `equal_equal`, `int`, `right_paren`, `right_brace`.

```monster
func example(int num) bool {
    for (int i = 0; i < num; i++) {
        print(num)
    }
    return (num == 3) 
}
```
