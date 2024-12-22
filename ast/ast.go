package ast

import "github.com/sean-d/sloth/token"

/*
Here we have three interfaces called Node, Statement and Expression.

Every node in our AST has to implement the Node interface, meaning it has to provide a TokenLiteral() method that
returns the literal value of the token it’s associated with.

TokenLiteral() will be used only for debugging and testing. The AST we are going to construct consists solely of Nodes
that are connected to each other - it’s a tree after all.

Some of these nodes implement the Statement and some the Expression interface. These interfaces only contain
dummy methods called statementNode and expressionNode respectively. They are not strictly necessary but help us by guiding
the Go compiler and possibly causing it to throw errors when we use a Statement where an Expression should’ve been used, and vice versa.

The Program node is going to be the root node of every AST our parser produces.
Every valid program is a series of statements. These statements are contained in the Program.Statements, which is just
a slice of AST nodes that implement the Statement interface.
*/

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

/*
Identifier is being used to hold the identifier of the binding, the x in let x = 5;. This implements the Expression interface.

But the identifier in a let statement doesn’t produce a value, right? So why is it an Expression?
It’s to keep things simple. Identifiers in other parts of a program do produce values, e.g.: let x = valueProducingIdentifier;.

And to keep the number of different node types small, we’ll use Identifier here to represent the name in a variable binding
and later reuse it, to represent an identifier as part of or as a complete expression.
*/
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// LetStatement has the fields we need: Name to hold the identifier of the binding and Value for the expression that produces the value.
// The two methods statementNode and TokenLiteral satisfy the Statement and Node interfaces respectively.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
