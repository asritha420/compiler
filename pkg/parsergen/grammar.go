package parsergen

type Grammar interface {
	generateFollowSet(ogByte byte, nt byte)
	generateFirstSet(ogByte byte, nt byte)
}
