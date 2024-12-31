package parser

import (
	"fmt"
	"github.com/sean-d/sloth/ast"
	"github.com/sean-d/sloth/lexer"
	"github.com/sean-d/sloth/token"
	"strconv"
)

// Setting the PEMDAS order of operations for later consideration.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // someFunction(X)
)

/*
Pratt Parser

A Pratt parser’s main idea is the association of parsing functions (which Pratt calls “semantic code”) with token types.
Whenever this token type is encountered, the parsing functions are called to parse the appropriate expression and
return an AST node that represents it.
Each token type can have up to two parsing functions associated with it, depending on whether the token is found in a prefix or an infix position.
*/

/*
Both of the following function types return an ast.Expression, since that’s what we’re here to parse.
Only the infixParseFn takes an argument: another ast.Expression. This argument is “left side” of the infix operator that’s being parsed.
A prefix operator doesn’t have a “left side”, per definition.

prefixParseFns gets called when we encounter the associated token type in prefix position and infixParseFn gets called
when we encounter the token type in infix position.
*/
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
)

/*
Parser has the following fields:
-lexer is a pointer to an instance of the lexer, on which we repeatedly call NextToken() to get the next token in the input.
-errors holds a slice of strings containing any errors the parsing encounters
-currentToken and peekToken act exactly like the two “pointers” our lexer has: position and readPosition.
-prefixParseFns and infixParseFns maps ensure the correct prefixParseFn or infixParseFn for the current token type

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

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New returns a pointer to a Parser
func New(l *lexer.Lexer) *Parser {
	parse := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// initialize the prefixParseFns map on Parser and register parsing functions:
	// EX: if we encounter a token of type token.IDENT the parsing function to call is parseIdentifier, a method we defined on *Parser.
	parse.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parse.registerPrefix(token.IDENT, parse.parseIdentifier)
	parse.registerPrefix(token.INT, parse.parseIntegerLiteral)
	parse.registerPrefix(token.BANG, parse.parsePrefixExpression)
	parse.registerPrefix(token.MINUS, parse.parsePrefixExpression)

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
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

// parseStatement checks the Type of the current token.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

/*
prefixParseFns gets called when we encounter the associated token type in prefix position and
infixParseFn gets called when we encounter the token type in infix position.
*/

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

/*
parseExpressionStatement builds an AST node and then attempts to fill its field by calling other parsing functions.
In this case there are a few differences though: we call parseExpression() with the constant LOWEST, and then we check
for an optional semicolon. Yes, it’s optional. If the peekToken is a token.SEMICOLON, we advance so it’s the curToken.
If it’s not there, that’s okay too, we don’t add an error to the parser if it’s not there.
Expression statements have optional semicolons (which makes it easier to type something like 5 + 5 into the REPL later on).
*/
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: nil,
	}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// parseExpression checks if a parsing function is associated with p.CurrentToken.Type in the prefix position.
// if true, the parsing function is called. if false, nil is returned.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	return leftExp
}

/*
parseIdentifier returns a *ast.Identifier with the current token in the Token field and the literal value of the token in Value.

Note: It doesn’t advance the tokens, it doesn’t call nextToken; we simply start with curToken being the type of token
you’re associated with and return with curToken being the last token that’s part of your expression type.
Never advance the tokens too far.
*/
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

// parseIntegerLiteral makes a call to strconv.ParseInt, which converts the string in p.curToken.Literal into an int64.
// The int64 then gets saved to the Value field, and we return the newly constructed *ast.IntegerLiteral node.
// If that doesn’t work, we add a new error to the parser’s errors field.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// noPrefixParseFnError just adds a formatted error message to our parser’s errors field.
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

/*
	parsePrefixExpression builds an AST node, in this case *ast.PrefixExpression, just like the parsing functions we saw before.

But then it does something different: it actually advances our tokens by calling p.nextToken().

When parsePrefixExpression is called, p.currentToken is either of type token.BANG or token.MINUS, because otherwise it
wouldn’t have been called. But in order to correctly parse a prefix expression like -5 more than one token has to be “consumed”.
So after using p.currentToken to build a *ast.PrefixExpression node, the method advances the tokens and calls parseExpression again.
This time with the precedence of prefix operators as argument.

Now, when parseExpression is called by parsePrefixExpression the tokens have been advanced and the current token is the
one after the prefix operator. In the case of -5, when parseExpression is called the p.currentToken.Type is token.INT.
parseExpression then checks the registered prefix parsing functions and finds parseIntegerLiteral, which builds
an *ast.IntegerLiteral node and returns it. parseExpression returns this newly constructed node and parsePrefixExpression
uses it to fill the Right field of *ast.PrefixExpression.
*/
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}
