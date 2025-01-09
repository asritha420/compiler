package aparsergen

type stack []rune

func (s *stack) push(val rune) {
	*s = append(*s, val)
}

func (s *stack) pop() rune {
	toPop := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return toPop
}
