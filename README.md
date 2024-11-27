# Compiler 

Compiler in Go 

## Standards
### Errors
Errors should at a MINIMUM specify the method they came from and give an error message that starts lowercase (unless the first thing needs to be capitalized). It should have the following structure:

`[source method] err msg`

If you would like to chain multiple errors together separate each one by a \n and use %w to print the error. For example:

`fmt.Errorf("[main] failed to convert rule because:\n%w", err)`

### Production Rules
#### Ranges
Note: In order to use the carrot `^` or `-` in a range, please escape it like so `\^`

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