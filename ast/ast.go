package ast

import (
	"bytes"
	"github.com/sean-d/sloth/token"
)

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
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

/*
ExpressionStatement has two fields: the Token field, which every node has, and the
Expression field, which holds the expression.

ast.ExpressionStatement fulfills the ast.Statement interface, which means we can add it to the Statements slice of ast.Program.
And that’s the whole reason why we’re adding ast.ExpressionStatement.
*/
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

// IntegerLiteral fulfills the ast.Expression interface, just like *ast.Identifier does, but there’s a notable difference
// to ast.Identifier in the structure itself: Value is an int64 and not a string. This is the field that’s going to
// contain the actual value the integer literal represents in the source code. When we build an *ast.IntegerLiteral
// we have to convert the string in *ast.IntegerLiteral.Token.Literal (which is something like "5") to an int64.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// PrefixExpression node has two noteworthy fields: Operator and Right. Operator is a string that’s going to contain
// either "-" or "!". The Right field contains the expression to the right of the operator.
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
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

/*
All String() for better debugging.

With these methods in place, we can now just call String() on *ast.Program and get our whole program back as a string.
That makes the structure of *ast.Program easily testable.
*/

// String creates a buffer and writes the return value of each statement’s String() method to it.
// And then it returns the buffer as a string. It delegates most of its work to the Statements of *ast.Program.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (i *Identifier) String() string { return i.Value }

// Deliberately add parentheses around the operator and its operand, the expression in Right.
// That allows us to see which operands belong to which operator.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
End String() business
*/

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

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
