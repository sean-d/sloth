package parser

import (
	"github.com/sean-d/sloth/ast"
	"github.com/sean-d/sloth/lexer"
	"github.com/sean-d/sloth/token"
)

/*
Parser has three fields: lexer, currentToken and peekToken.
-lexer is a pointer to an instance of the lexer, on which we repeatedly call NextToken() to get the next token in the input.
-currentToken and peekToken act exactly like the two “pointers” our lexer has: position and readPosition.

Instead of pointing to a character in the input, they point to the current and the next token.

Both are important: we need to look at the currentToken, which is the current token under examination,
to decide what to do next, and we also need peekToken for this decision if currentToken doesn’t give us enough information.

Think of a single line only containing 5;. Then currentToken is a token.INT and we need peekToken to decide whether
we are at the end of the line or if we are at just the start of an arithmetic expression.
*/
type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	parse := &Parser{
		lexer: l,
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
	default:
		return nil
	}
}

// TODO start at bottom of page 40 with parseLetStatement()
