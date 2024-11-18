package parsergen

type Grammar interface {
	generateFollowSet(ogByte byte, nt byte)
	generateFirstSet(ogByte byte, nt byte)
}

type Rule struct {
	//TODO: everywhere I use a rule in code, make sure I call it name or production for naming clarity everywhere
	name       byte //this cant support epsilon
	production []byte
}

// TODO: can implement some checking in here, return error in invalid rule provided
func NewRule(name byte, production []byte) *Rule {

	return &Rule{
		name:       name,
		production: production,
	}

}
