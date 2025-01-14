package regex

type Regex string

func (r *Regex) GetAST() {

}

type regexParser struct {
	regex []rune
	curr  int
}

func Parse(regex []rune) {

}
