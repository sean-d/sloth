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
	INT   = "INT"   // 0123456789

	//operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	//delimeters
	COMMA     = ","
	SEMICOLON = ";"

	//groupings
	QUOTES = "\""
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks the keywords table to see if a given identifier is a keyword.
// If so, the TokeType of that keyword is returned. If not, token.IDENT is returned which is the
// TokenType for all user-defined identifiers
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
