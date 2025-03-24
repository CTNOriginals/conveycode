package lexer

type StateFn func(*lexer) StateFn

func LexText(lx *lexer) StateFn {

	return nil
}
