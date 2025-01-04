package ast

import (
	"bytes"
	"github.com/sean-d/sloth/token"
	"strings"
)

/*
All String() for better debugging.

With these methods in place, we can now just call String() on *ast.Program and get our whole program back as a string.
That makes the structure of *ast.Program easily testable.

String creates a buffer and writes the return value of each statement’s String() method to it.
And then it returns the buffer as a string. It delegates most of its work to the Statements of *ast.Program.
*/

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

// Program section
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Let Statement section

// LetStatement has the fields we need: Name to hold the identifier of the binding and Value for the expression that produces the value.
// The two methods statementNode and TokenLiteral satisfy the Statement and Node interfaces respectively.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
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

// implicit stuff for implicit things
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Return statement section
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
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

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// Expression statement stuff

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

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// Block statement stuff

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// Identifier expression stuff

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

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// Boolean expression stuff
// Boolean: The Value field can hold values of the type bool, which means that we’re going to save
// either true or false in there (the Go bool values, not the Sloth literals).
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// Integer literal stuff

// IntegerLiteral fulfills the ast.Expression interface, just like *ast.Identifier does, but there’s a notable difference
// to ast.Identifier in the structure itself: Value is an int64 and not a string. This is the field that’s going to
// contain the actual value the integer literal represents in the source code. When we build an *ast.IntegerLiteral
// we have to convert the string in *ast.IntegerLiteral.Token.Literal (which is something like "5") to an int64.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// StringLiteral fulfills the ast.Expression interface, just like *ast.Identifier does
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// PrefixExpression stuff
// PrefixExpression node has two noteworthy fields: Operator and Right. Operator is a string that’s going to contain
// either "-" or "!". The Right field contains the expression to the right of the operator.
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

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

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// InfixExpression stuff
// InfixExpression fulfills the ast.Expression and ast.Node interfaces, by defining the expressionNode(), TokenLiteral() and String() methods.
// The only difference to ast.PrefixExpression is the new field called Left, which can hold any expression.
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// IfExpression fulfills the ast.Expression interface and has three fields that can represent an if-else-conditional.
// Condition holds the condition, which can be any expression, and Consequence and Alternative point to the consequence
// and alternative of the conditional respectively.
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// Function literal stuff
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// CallExpression consists of an expression that results in a function when evaluated and a list of expressions
// that are the arguments to this function call.
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
