package lexer

import (
	"github.com/sean-d/sloth/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a pointer to a Lexer that is instantiated with the possible inputs
// The new Lexer has an input with the rest being 0.
// readChar() is called to have ch represent the first char in the Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

// NextToken works as follows:
// We look at the current character under examination (l.ch) and return a token depending on which character it is.
// Before returning the token we advance our pointers into the input so when we call NextToken() again the l.ch field is already updated.
//
// Default branch captures non-recognized characters as defined in token.go.
// We identify if a character is a letter and if so, it needs to keep reading until a non-letter occurs. This signifies the end
// of a keyword or identifier and then we sort our if what was just read is a keyword or identifier so the correct token type is used.
//
// We early exit in default: when calling readIdentifier/readNumber, we call readChar repeatedly and advances the readPosition/position fields
// beyond the last character of the current identifier. Because of this, we don't need to call readChar() after the switch again.
// If we wind up at the token.ILLEGAL we have something we have no idea what to do with.
//
// There are special cases with double-characters. As more 2-char literals become special, We will eventually use a helper function to determine if the char following a specific
// char is what we desire and if so, we build a string from them both and return a token.Token with that string and associated
// token type.
//
// A small function called newToken helps us with initializing these tokens.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

// skipWhitespace will determine if the current character is a space, a newline, a tab, or a return
// and call readChar to get the next character.
//
// Whitespace is skipped/ignored rather than used. EMBRACE THE CHAOS
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readChar provides the next character and advances the position in the input string.
// 1. checks if the end of input has been reached
// 1a. if so, l.ch gets set to 0 and signals nothing has been read or EOF
// 1b. if EOF is not true, l.ch gets set to the next char by accessing l.input[l.readPosition]
//
// 2. l.position is updated to the just used l.readPosition and l.readPosition is incremented by one.
// This way, l.readPosition will always point to the next position that will be read from
// and l.position always points to the position last read.
//
// We are only supporting ASCII to keep thing simple
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar is really similar to readChar, except that it doesn’t increment l.position and l.readPosition.
// We only want to “peek” ahead in the input and not move around in it, so we know what a call to readChar() would return.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// readIdentifier reads in an identifier and advances the lexer position until it encounters a non-letter character
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber only takes in ints. we are not worrying about any other numbers. who cares :)
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString calls readChar until it encounters either a closing double quote or the end of the input.
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// isLetter returns true if the passed in character is a->z or A->Z or is a underscore.
// we allow underscores so we can snake_case things :)
func isLetter(ch byte) bool {
	return 'a' <= ch && 'z' >= ch || 'A' <= ch && 'Z' >= ch || ch == '_'
}

// isDigit returns true if the passed in character is greater than 0 and less than 9
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newToken is a helper function that takes in a token type and the literal
// and returns the token for that
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
