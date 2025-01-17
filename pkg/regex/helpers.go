package regex

const epsilon = rune(0)

// TODO: fix these literals, so they match w/ the function name names in grammar, inconsistent naming?
var (
	anyChar   = append(append(digits(), lowercaseLetters()...), uppercaseLetters()...)
	firstSets = map[string][]rune{
		"Regex":         append(anyChar, '(', '['),
		"Alt":           append(anyChar, '(', '['),
		"AltPrime":      {'|', epsilon},
		"Concat":        append(anyChar, '(', '['),
		"ConcatPrime":   append(anyChar, '(', '[', epsilon),
		"Repeat":        append(anyChar, '(', '['),
		"Quantifier":    {'*', '+', '?'},
		"Group":         append(anyChar, '(', '['),
		"CharRange":     {'['},
		"CharRangeBody": append(anyChar, '^'),
		"CharRangeAtom": append(anyChar),
		"Char":          append(anyChar),
	}
)

// TODO: make it actually unicode support ranges
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
