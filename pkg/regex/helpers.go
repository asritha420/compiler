package regex

// TODO: fix these literals, so they match w/ the function name names in grammar, inconsistent naming?
var (
	anyChar   = append(append(digits(), lowercaseLetters()...), uppercaseLetters()...)
	firstSets = map[string][]rune{
		"Regex":         append(anyChar, '(', '['),
		"Alt":           append(anyChar, '(', '['),
		"AltPrime":      {'|'},
		"Concat":        append(anyChar, '(', '['),
		"ConcatPrime":   append(anyChar, '(', '['),
		"Repeat":        append(anyChar, '(', '['),
		"Quantifier":    {'*', '+', '?'},
		"Group":         append(anyChar, '(', '['),
		"CharRange":     {'['},
		"CharRangeBody": append(anyChar, '^'),
		"CharRangeAtom": append(anyChar, '(', '['),
		"Char":          append(anyChar, '(', '['),
	}
)

func digits() []rune {
	digits := make([]rune, 10)

	for i := 0; i < 10; i++ {
		digits[i] = rune(i)
	}

	return digits
}

func lowercaseLetters() []rune {
	lowercaseLetters := make([]rune, 26)

	for i := 0; i < 26; i++ {
		lowercaseLetters[i] = 'a' + rune(i)
	}

	return lowercaseLetters
}

func uppercaseLetters() []rune {
	uppercaseLetters := make([]rune, 26)

	for i := 0; i < 26; i++ {
		uppercaseLetters[i] = 'A' + rune(i)
	}

	return uppercaseLetters
}
