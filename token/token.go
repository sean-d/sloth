package token

type TokenType string

// Token holds:
// - the type of token: integer, right-bracket
// - the literal value of the token: 5, ]
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// ILLEGAL signifies a token/char we don't know about
	// EOF stands for end of file and lets the parser know when to stop
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//identifiers + literals
	IDENT = "IDENT" //add, someName, x, y...
	INT   = "INT"   // 12345

	//operators
	ASSIGN = "="
	PLUS   = "+"

	//delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
