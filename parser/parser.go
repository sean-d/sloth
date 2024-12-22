package parser

import (
	"fmt"
	"github.com/sean-d/sloth/ast"
	"github.com/sean-d/sloth/lexer"
	"github.com/sean-d/sloth/token"
)

/*
Parser has the following fields: lexer, errors, currentToken and peekToken.
-lexer is a pointer to an instance of the lexer, on which we repeatedly call NextToken() to get the next token in the input.
-errors holds a slice of strings containing any errors the parsing encounters
-currentToken and peekToken act exactly like the two “pointers” our lexer has: position and readPosition.

Instead of pointing to a character in the input, they point to the current and the next token.

Both are important: we need to look at the currentToken, which is the current token under examination,
to decide what to do next, and we also need peekToken for this decision if currentToken doesn’t give us enough information.

Think of a single line only containing 5;. Then currentToken is a token.INT and we need peekToken to decide whether
we are at the end of the line or if we are at just the start of an arithmetic expression.
*/
type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

// New returns a pointer to a Parser
func New(l *lexer.Lexer) *Parser {
	parse := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// Read two tokens to set both currentToken and peekToken
	parse.nextToken()
	parse.nextToken()

	return parse
}

// nextToken is a small helper that advances both currentToken and peekToken
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

/*
ParseProgram constructs the root node of the AST, an *ast.Program. It then iterates over every token in the input until
it encounters a token.EOF token. It does this by repeatedly calling nextToken, which advances both p.curToken and p.peekToken.
In every iteration it calls parseStatement, whose job it is to parse a statement. If parseStatement returned something
other than nil, an ast.Statement, its return value is added to Statements slice of the AST root node.
When nothing is left to parse the *ast.Program root node is returned.
*/
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement checks the Type of the current token. If the currentToken is a "let", parseLetStatement is called.
// Nil is returned otherwise.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

/*
parseLetStatement constructs an *ast.LetStatement node with the token it’s currently sitting on (a token.LET token) and
then advances the tokens while making assertions about the next token with calls to expectPeek.

First it expects a token.IDENT token, which it then uses to construct an *ast.Identifier node. Then it expects an
equal sign, and finally it jumps over the expression following the equal sign until it encounters a semicolon.

The skipping of expressions will be replaced, of course, as soon as we know how to parse them.
*/
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: expressions are being skipped until a semicolon is encountered
	// also update the above function doc string

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

}

// parseReturnStatement constructs an ast.ReturnStatement, with the current token it’s sitting on as Token.
// It then brings the parser in place for the expression that comes next by calling nextToken() and finally,
// there’s the cop-out. It skips over every expression until it encounters a semicolon. That’s it.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.currentToken}
	p.nextToken()
	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// currentTokenIs returns the bool repr of asserting if the current token is of an assumed type
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs returns the bool repr of asserting if the next token is of an assumed type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

/*
expectPeek method is one of the “assertion functions” nearly all parsers share. Their primary purpose is to enforce
the correctness of the order of tokens by checking the type of the next token.

Our expectPeek here checks the type of the peekToken and only if the type is correct does it advance the tokens by
calling nextToken.
*/

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Errors returns a slice of strings containing all parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError adds an error to p.errors when the type of peekToken does not match the expectation.
func (p *Parser) peekError(tok token.TokenType) {
	message := fmt.Sprintf("expected next token to be %s, got %s instead", tok, p.peekToken.Type)

	p.errors = append(p.errors, message)
}
