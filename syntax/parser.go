package syntax

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const trace = false

// We represent x++, x-- as assignments x += ImplicitOne, x -= ImplicitOne.
// ImplicitOne should not be used elsewhere.
var ImplicitOne = &BasicLit{Value: "1"}

// Mode describes the parser mode.
type Mode uint

type FilenameHandler func(name string) string

// The stopset contains keywords that start a statement.
// They are good synchronization points in case of syntax
// errors and (usually) shouldn't be skipped over.
const stopset uint64 = 1<<_Break |
	1<<_Const |
	1<<_Continue |
	1<<_Defer |
	1<<_Fallthrough |
	1<<_For |
	1<<_Go |
	1<<_Goto |
	1<<_If |
	1<<_Return |
	1<<_Select |
	1<<_Switch |
	1<<_Type |
	1<<_Var

// Modes supported by the parser.
const (
	CheckBranches Mode = 1 << iota // check correct use of labels, break, continue, and goto statements
)

type parser struct {
	file *PosBase
	errh ErrorHandler
	mode Mode
	scanner

	base   *PosBase // current position base
	first  error    // first error encountered
	errcnt int      // number of errors encountered
	pragma Pragma   // pragma flags

	fnest  int    // function nesting level (for error handling)
	xnest  int    // expression nesting level (for complit ambiguity resolution)
	indent []byte // tracing support
}

func (p *parser) init(file *PosBase, r io.Reader, errh ErrorHandler) {
	p.file = file
	p.errh = errh
	// p.mode = mode
	p.scanner.init(
		r,
		// Error and directive handler for scanner.
		// Because the (line, col) positions passed to the
		// handler is always at or after the current reading
		// position, it is safe to use the most recent position
		// base to compute the corresponding Pos value.
		func(line, col uint, msg string) {
			if msg[0] != '/' {
				p.errorAt(p.posAt(line, col), msg)
				return
			}

			// otherwise it must be a comment containing a line or go: directive
			text := commentText(msg)
			if strings.HasPrefix(text, "line ") {
				var pos Pos // position immediately following the comment
				if msg[1] == '/' {
					// line comment (newline is part of the comment)
					pos = MakePos(p.file, line+1, colbase)
				} else {
					// regular comment
					// (if the comment spans multiple lines it's not
					// a valid line directive and will be discarded
					// by updateBase)
					pos = MakePos(p.file, line, col+uint(len(msg)))
				}
				p.updateBase(pos, line, col+2+5, text[5:]) // +2 to skip over // or /*
				return
			}

			// go: directive (but be conservative and test)
			// if pragh != nil && strings.HasPrefix(text, "go:") {
			// 	p.pragma |= pragh(p.posAt(line, col+2), text) // +2 to skip over // or /*
			// }
		},
		directives,
	)

	p.base = file
	p.first = nil
	p.errcnt = 0
	p.pragma = 0

	p.fnest = 0
	p.xnest = 0
	p.indent = nil
}

//Parse parse file to AST
func Parse(base *PosBase, src io.Reader, errh ErrorHandler) (_ *File, first error) {
	defer func() {
		if p := recover(); p != nil {
			// if err, ok := p.(Error); ok {
			// 	first = err
			// 	return
			// }
			panic(p)
		}
	}()

	var p parser
	p.init(base, src, errh)
	p.next()
	return p.fileOrNil(), p.first
}

// ----------------------------------------------------------------------------
// Package files
//
// Parse methods are annotated with matching Go productions as appropriate.
// The annotations are intended as guidelines only since a single Go grammar
// rule may be covered by multiple parse methods and vice versa.
//
// Excluding methods returning slices, parse methods named xOrNil may return
// nil; all others are expected to return a valid non-nil node.

// SourceFile = PackageClause ";" { ImportDecl ";" } { TopLevelDecl ";" } .
func (p *parser) fileOrNil() *File {
	if trace {
		defer p.trace("file")()
	}

	f := new(File)
	f.pos = p.pos()

	// PackageClause
	if !p.got(_Package) {
		p.syntaxError("package statement must be first")
		return nil
	}
	f.PkgName = p.name()
	p.want(_Semi)

	// don't bother continuing if package clause has errors
	if p.first != nil {
		return nil
	}

	// { ImportDecl ";" }
	for p.got(_Import) {
		f.DeclList = p.appendGroup(f.DeclList, p.importDecl)
		p.want(_Semi)
	}

	l := []Stmt{}
	// { TopLevelDecl ";" }
	for p.tok != _EOF {
		switch p.tok {
		// case _Var:
		// 	p.next()
		// 	f.DeclList = p.appendGroup(f.DeclList, p.varDecl)

		case _Func:
			p.next()
			if d := p.funcDeclOrNil(); d != nil {
				f.DeclList = append(f.DeclList, d)
				addFunc(d)
			}

		default:
			//TODO   define the block
			// if p.tok == _Lbrace && len(f.DeclList) > 0 && isEmptyFuncDecl(f.DeclList[len(f.DeclList)-1]) {
			// 	// opening { of function declaration on next line
			// 	// p.syntaxError("unexpected semicolon or newline before {")
			// } else {
			// 	// p.syntaxError("non-declaration statement outside function body")
			// }
			s := p.stmtOrNil()
			if s == nil {
				break
			}
			l = append(l, s)
			// ";" is optional before "}"
			if !p.got(_Semi) && p.tok != _Rbrace {
				// p.syntaxError("at end of statement")
				p.advance(_Semi, _Rbrace, _Case, _Default)
				p.got(_Semi) // avoid spurious empty statement
			}
			// p.advance(_Const, _Type, _Var, _Func)
			continue
		}

		// Reset p.pragma BEFORE advancing to the next token (consuming ';')
		// since comments before may set pragmas for the next function decl.
		p.pragma = 0

		// if p.tok != _EOF && !p.got(_Semi) {
		// 	// p.syntaxError("after top level declaration")
		// 	p.advance(_Const, _Type, _Var, _Func)
		// }
	}

	f.Lines = p.source.line
	f.Blocks = l
	return f
}

func isEmptyFuncDecl(dcl Decl) bool {
	f, ok := dcl.(*FuncDecl)
	return ok && f.Body == nil
}

// FunctionDecl = "func" FunctionName ( Function | Signature ) .
// FunctionName = identifier .
// Function     = Signature FunctionBody .
// MethodDecl   = "func" Receiver MethodName ( Function | Signature ) .
// Receiver     = Parameters .
func (p *parser) funcDeclOrNil() *FuncDecl {
	if trace {
		defer p.trace("funcDecl")()
	}

	f := new(FuncDecl)
	f.pos = p.pos()
	fmt.Println(p.lit)
	// if p.tok == _Lparen {
	// 	rcvr := p.paramList()
	// 	switch len(rcvr) {
	// 	case 0:
	// 		p.error("method has no receiver")
	// 	default:
	// 		p.error("method has multiple receivers")
	// 		fallthrough
	// 	case 1:
	// 		f.Recv = rcvr[0]
	// 	}
	// }

	if p.tok != _Name {
		p.syntaxError("expecting func name ")
		p.advance(_Lbrace, _Semi)
		return nil
	}

	f.Name = p.name()
	f.Type = p.funcType()
	if p.tok == _Lbrace {
		f.Body = p.funcBody()
	}
	f.Pragma = p.pragma

	return f
}

// Parameters    = "(" [ ParameterList [ "," ] ] ")" .
// ParameterList = ParameterDecl { "," ParameterDecl } .
func (p *parser) paramList() (list []*Field) {
	if trace {
		defer p.trace("paramList")()
	}

	// pos := p.pos()

	var named int // number of parameters that have an explicit name and type
	p.list(_Lparen, _Comma, _Rparen, func() bool {
		if par := p.paramDeclOrNil(); par != nil {
			if par.Name == nil {
				panic("parameter without name ")
			}
			if par.Name != nil {
				named++
			}
			list = append(list, par)
		}
		return false
	})

	// distribute parameter types
	// if named == 0 {
	// 	// all unnamed => found names are named types
	// 	for _, par := range list {
	// 		if typ := par.Name; typ != nil {
	// 			par.Type = typ
	// 			par.Name = nil
	// 		}
	// 	}
	// } else if named != len(list) {
	// 	// some named => all must be named
	// 	ok := true
	// 	var typ Expr
	// 	for i := len(list) - 1; i >= 0; i-- {
	// 		if par := list[i]; par.Type != nil {
	// 			typ = par.Type
	// 			if par.Name == nil {
	// 				ok = false
	// 				n := p.newName("_")
	// 				n.pos = typ.Pos() // correct position
	// 				par.Name = n
	// 			}
	// 		} else if typ != nil {
	// 			par.Type = typ
	// 		} else {
	// 			// par.Type == nil && typ == nil => we only have a par.Name
	// 			ok = false
	// 			t := p.bad()
	// 			t.pos = par.Name.Pos() // correct position
	// 			par.Type = t
	// 		}
	// 	}
	// 	if !ok {
	// 		// p.syntaxError_at(pos, "mixed named and unnamed function parameters")
	// 	}
	// }

	return
}

// ParameterDecl = [ IdentifierList ] [ "..." ] Type .
func (p *parser) paramDeclOrNil() *Field {
	if trace {
		defer p.trace("paramDecl")()
	}

	f := new(Field)
	f.pos = p.pos()
	switch p.tok {
	case _Name:
		f.Name = p.name()
		p.next()
		// switch p.tok {
		// case _Name, _Star, _Arrow, _Func, _Lbrack, _Chan, _Map, _Struct, _Interface, _Lparen:
		// 	// sym name_or_type
		// 	f.Type = p.type_()

		// case _Dot:
		// 	// name_or_type
		// 	// from dotname
		// 	f.Type = p.dotname(f.Name)
		// 	f.Name = nil
		// }

	// case _Arrow, _Star, _Func, _Lbrack, _Chan, _Map, _Struct, _Interface, _Lparen:
	// 	// name_or_type
	// 	f.Type = p.type_()

	default:
		p.syntaxError("expecting )")
		p.advance(_Comma, _Rparen)
		return nil
	}

	return f
}

func (p *parser) want(tok token) {
	if !p.got(tok) {
		// p.syntaxError("expecting " + tokstring(tok))
		p.advance()
	}
}

// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func (p *parser) varDecl(group *Group) Decl {
	if trace {
		defer p.trace("varDecl")()
	}

	d := new(VarDecl)
	d.pos = p.pos()

	d.NameList = p.nameList(p.name())
	if p.got(_Assign) {
		d.Values = p.exprList()
	} else {
		d.Type = p.type_()
		if p.got(_Assign) {
			d.Values = p.exprList()
		}
	}
	d.Group = group

	return d
}

func (p *parser) type_() Expr {
	if trace {
		defer p.trace("type_")()
	}

	typ := p.typeOrNil()
	if typ == nil {
		typ = p.bad()
		p.syntaxError("expecting type")
		p.advance()
	}

	return typ
}

// typeOrNil is like type_ but it returns nil if there was no type
// instead of reporting an error.
//
// Type     = TypeName | TypeLit | "(" Type ")" .
// TypeName = identifier | QualifiedIdent .
// TypeLit  = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
// 	      SliceType | MapType | Channel_Type .
func (p *parser) typeOrNil() Expr {
	if trace {
		defer p.trace("typeOrNil")()
	}

	pos := p.pos()
	switch p.tok {

	case _Func:
		// fntype
		p.next()
		return p.funcType()

	case _Lbrack:
		// '[' oexpr ']' ntype
		// '[' _DotDotDot ']' ntype
		p.next()
		p.xnest++
		if p.got(_Rbrack) {
			// []T
			p.xnest--
			t := new(SliceType)
			t.pos = pos
			t.Elem = p.type_()
			return t
		}

		// [n]T
		t := new(ArrayType)
		t.pos = pos
		if !p.got(_DotDotDot) {
			t.Len = p.expr()
		}
		p.want(_Rbrack)
		p.xnest--
		t.Elem = p.type_()
		return t

	case _Map:
		// _Map '[' ntype ']' ntype
		p.next()
		p.want(_Lbrack)
		t := new(MapType)
		t.pos = pos
		t.Key = p.type_()
		p.want(_Rbrack)
		t.Value = p.type_()
		return t

	case _Name:
		return p.dotname(p.name())

	case _Lparen:
		p.next()
		t := p.type_()
		p.want(_Rparen)
		return t
	}

	return nil
}

func (p *parser) dotname(name *Name) Expr {
	if trace {
		defer p.trace("dotname")()
	}

	if p.tok == _Dot {
		s := new(SelectorExpr)
		s.pos = p.pos()
		p.next()
		s.X = name
		s.Sel = p.name()
		return s
	}
	return name
}

// ExpressionList = Expression { "," Expression } .
func (p *parser) exprList() Expr {
	if trace {
		defer p.trace("exprList")()
	}

	x := p.expr()
	if p.got(_Comma) {
		list := []Expr{x, p.expr()}
		for p.got(_Comma) {
			list = append(list, p.expr())
		}
		t := new(ListExpr)
		t.pos = x.Pos()
		t.ElemList = list
		x = t
	}
	return x
}

func (p *parser) expr() Expr {
	if trace {
		defer p.trace("expr")()
	}

	return p.binaryExpr(0)
}

// Expression = UnaryExpr | Expression binary_op Expression .
func (p *parser) binaryExpr(prec int) Expr {
	// don't trace binaryExpr - only leads to overly nested trace output

	x := p.unaryExpr()
	for (p.tok == _Operator || p.tok == _Star) && p.prec > prec {
		t := new(Operation)
		t.pos = p.pos()
		t.Op = p.op
		t.X = x
		tprec := p.prec
		p.next()
		t.Y = p.binaryExpr(tprec)
		x = t
	}
	return x
}

// UnaryExpr = PrimaryExpr | unary_op UnaryExpr .
func (p *parser) unaryExpr() Expr {
	if trace {
		defer p.trace("unaryExpr")()
	}

	switch p.tok {
	case _Operator, _Star:
		switch p.op {
		case Mul, Add, Sub, Not, Xor:
			x := new(Operation)
			x.pos = p.pos()
			x.Op = p.op
			p.next()
			x.X = p.unaryExpr()
			return x

		case And:
			x := new(Operation)
			x.pos = p.pos()
			x.Op = And
			p.next()
			// unaryExpr may have returned a parenthesized composite literal
			// (see comment in operand) - remove parentheses if any
			x.X = unparen(p.unaryExpr())
			return x
		}

	case _Arrow:
		// receive op (<-x) or receive-only channel (<-chan E)
		pos := p.pos()
		p.next()

		// If the next token is _Chan we still don't know if it is
		// a channel (<-chan int) or a receive op (<-chan int(ch)).
		// We only know once we have found the end of the unaryExpr.

		x := p.unaryExpr()

		// There are two cases:
		//
		//   <-chan...  => <-x is a channel type
		//   <-x        => <-x is a receive operation
		//
		// In the first case, <- must be re-associated with
		// the channel type parsed already:
		//
		//   <-(chan E)   =>  (<-chan E)
		//   <-(chan<-E)  =>  (<-chan (<-E))

		// x is not a channel type => we have a receive op
		o := new(Operation)
		o.pos = pos
		o.Op = Recv
		o.X = x
		return o
	}

	// TODO(mdempsky): We need parens here so we can report an
	// error for "(x) := true". It should be possible to detect
	// and reject that more efficiently though.
	return p.pexpr(true)
}

func (p *parser) pexpr(keep_parens bool) Expr {
	if trace {
		defer p.trace("pexpr")()
	}

	x := p.operand(keep_parens)

loop:
	for {
		pos := p.pos()
		switch p.tok {
		case _Dot:
			p.next()
			switch p.tok {
			case _Name:
				// pexpr '.' sym
				t := new(SelectorExpr)
				t.pos = pos
				t.X = x
				t.Sel = p.name()
				x = t

			case _Lparen:
				t := new(AssertExpr)
				t.pos = pos
				t.X = x
				t.Type = p.expr()
				x = t
				p.want(_Rparen)

			default:
				p.syntaxError("expecting name or (")
				p.advance(_Semi, _Rparen)
			}

		case _Lbrack:
			p.next()
			p.xnest++

			var i Expr
			if p.tok != _Colon {
				i = p.expr()
				if p.got(_Rbrack) {
					// x[i]
					t := new(IndexExpr)
					t.pos = pos
					t.X = x
					t.Index = i
					x = t
					p.xnest--
					break
				}
			}

			// x[i:...
			t := new(SliceExpr)
			t.pos = pos
			t.X = x
			t.Index[0] = i
			p.want(_Colon)
			if p.tok != _Colon && p.tok != _Rbrack {
				// x[i:j...
				t.Index[1] = p.expr()
			}
			if p.got(_Colon) {
				t.Full = true
				// x[i:j:...]
				if t.Index[1] == nil {
					p.error("middle index required in 3-index slice")
				}
				if p.tok != _Rbrack {
					// x[i:j:k...
					t.Index[2] = p.expr()
				} else {
					p.error("final index required in 3-index slice")
				}
			}
			p.want(_Rbrack)

			x = t
			p.xnest--

		case _Lparen:
			t := new(CallExpr)
			t.pos = pos
			t.Fun = x
			t.ArgList, t.HasDots = p.argList()
			x = t

		case _Lbrace:
			// operand may have returned a parenthesized complit
			// type; accept it but complain if we have a complit
			t := unparen(x)
			// determine if '{' belongs to a composite literal or a block statement
			complit_ok := false
			switch t.(type) {
			case *Name, *SelectorExpr:
				if p.xnest >= 0 {
					// x is considered a composite literal type
					complit_ok = true
				}
			case *ArrayType:
				// x is a comptype
				complit_ok = true
			}
			if !complit_ok {
				break loop
			}
			if t != x {
				// p.syntaxError("cannot parenthesize type in composite literal")
				// already progressed, no need to advance
			}
			n := p.complitexpr()
			n.Type = x
			x = n

		default:
			break loop
		}
	}

	return x
}

// LiteralValue = "{" [ ElementList [ "," ] ] "}" .
func (p *parser) complitexpr() *CompositeLit {
	if trace {
		defer p.trace("complitexpr")()
	}

	x := new(CompositeLit)
	x.pos = p.pos()

	p.xnest++
	x.Rbrace = p.list(_Lbrace, _Comma, _Rbrace, func() bool {
		// value
		e := p.bare_complitexpr()
		if p.tok == _Colon {
			// key ':' value
			l := new(KeyValueExpr)
			l.pos = p.pos()
			p.next()
			l.Key = e
			l.Value = p.bare_complitexpr()
			e = l
			x.NKeys++
		}
		x.ElemList = append(x.ElemList, e)
		return false
	})
	p.xnest--

	return x
}

// Element = Expression | LiteralValue .
func (p *parser) bare_complitexpr() Expr {
	if trace {
		defer p.trace("bare_complitexpr")()
	}

	if p.tok == _Lbrace {
		// '{' start_complit braced_keyval_list '}'
		return p.complitexpr()
	}

	return p.expr()
}

func (p *parser) argList() (list []Expr, hasDots bool) {
	if trace {
		defer p.trace("argList")()
	}

	p.xnest++
	p.list(_Lparen, _Comma, _Rparen, func() bool {
		list = append(list, p.expr())
		hasDots = p.got(_DotDotDot)
		return hasDots
	})
	p.xnest--

	return
}

func (p *parser) oliteral() *BasicLit {
	if p.tok == _Literal {
		b := new(BasicLit)
		b.pos = p.pos()
		b.Value = p.lit
		b.Kind = p.kind
		p.next()
		return b
	}
	return nil
}

// Operand     = Literal | OperandName | MethodExpr | "(" Expression ")" .
// Literal     = BasicLit | CompositeLit | FunctionLit .
// BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
// OperandName = identifier | QualifiedIdent.
func (p *parser) operand(keep_parens bool) Expr {
	if trace {
		defer p.trace("operand " + p.tok.String())()
	}

	switch p.tok {
	case _Name:
		return p.name()

	case _Literal:
		return p.oliteral()

	case _Lparen:
		pos := p.pos()
		p.next()
		p.xnest++
		x := p.expr()
		p.xnest--
		p.want(_Rparen)

		// Optimization: Record presence of ()'s only where needed
		// for error reporting. Don't bother in other cases; it is
		// just a waste of memory and time.

		// Parentheses are not permitted on lhs of := .
		// switch x.Op {
		// case ONAME, ONONAME, OPACK, OTYPE, OLITERAL, OTYPESW:
		// 	keep_parens = true
		// }

		// Parentheses are not permitted around T in a composite
		// literal T{}. If the next token is a {, assume x is a
		// composite literal type T (it may not be, { could be
		// the opening brace of a block, but we don't know yet).
		if p.tok == _Lbrace {
			keep_parens = true
		}

		// Parentheses are also not permitted around the expression
		// in a go/defer statement. In that case, operand is called
		// with keep_parens set.
		if keep_parens {
			px := new(ParenExpr)
			px.pos = pos
			px.X = x
			x = px
		}
		return x

	case _Func:
		pos := p.pos()
		p.next()
		t := p.funcType()
		if p.tok == _Lbrace {
			p.xnest++

			f := new(FuncLit)
			f.pos = pos
			f.Type = t
			f.Body = p.funcBody()

			p.xnest--
			return f
		}
		return t

	default:
		x := p.bad()
		// p.syntaxError("expecting expression")
		p.advance()
		return x
	}

	// Syntactically, composite literals are operands. Because a complit
	// type may be a qualified identifier which is handled by pexpr
	// (together with selector expressions), complits are parsed there
	// as well (operand is only called from pexpr).
}

func (p *parser) bad() *BadExpr {
	b := new(BadExpr)
	b.pos = p.pos()
	return b
}

// context must be a non-empty string unless we know that p.tok == _Lbrace.
func (p *parser) blockStmt(context string) *BlockStmt {
	if trace {
		defer p.trace("blockStmt")()
	}

	s := new(BlockStmt)
	s.pos = p.pos()

	// people coming from C may forget that braces are mandatory in Go
	if !p.got(_Lbrace) {
		// p.syntaxError("expecting { after " + context)
		p.advance(_Name, _Rbrace)
		s.Rbrace = p.pos() // in case we found "}"
		if p.got(_Rbrace) {
			return s
		}
	}

	s.List = p.stmtList()
	s.Rbrace = p.pos()
	p.want(_Rbrace)

	return s
}

// main package block, ignore funcstmt
func (p *parser) mainBlockStmt(context string) *BlockStmt {
	if trace {
		defer p.trace("mainBlockStmt")()
	}

	s := new(BlockStmt)
	s.pos = p.pos()

	// people coming from C may forget that braces are mandatory in Go
	if !p.got(_Lbrace) {
		// p.syntaxError("expecting { after " + context)
		p.advance(_Name, _Rbrace)
		s.Rbrace = p.pos() // in case we found "}"
		if p.got(_Rbrace) {
			return s
		}
	}

	s.List = p.stmtList()
	s.Rbrace = p.pos()
	p.want(_Rbrace)

	return s
}

// StatementList = { Statement ";" } .
func (p *parser) stmtList() (l []Stmt) {
	if trace {
		defer p.trace("stmtList")()
	}

	for p.tok != _EOF && p.tok != _Rbrace && p.tok != _Case && p.tok != _Default {
		s := p.stmtOrNil()
		if s == nil {
			break
		}
		l = append(l, s)
		// ";" is optional before "}"
		if !p.got(_Semi) && p.tok != _Rbrace && !p.isCondition(s) {
			// p.syntaxError("at end of statement")
			p.advance(_Semi, _Rbrace, _Case, _Default)
			p.got(_Semi) // avoid spurious empty statement
		}
	}
	return
}

func (p *parser) isCondition(s Stmt) bool {
	if _, ok := s.(*IfStmt); ok {
		return true
	}
	return false
}

func (p *parser) declStmt(f func(*Group) Decl) *DeclStmt {
	if trace {
		defer p.trace("declStmt")()
	}

	s := new(DeclStmt)
	s.pos = p.pos()

	p.next() // _Const, _Type, or _Var
	s.Decl = p.appendGroup(nil, f)[0]

	return s
}

func (p *parser) stmtOrNil() Stmt {
	if trace {
		defer p.trace("stmt " + p.tok.String())()
	}

	// // Most statements (assignments) start with an identifier;
	// // look for it first before doing anything more expensive.
	if p.tok == _Name {
		lhs := p.exprList()
		if label, ok := lhs.(*Name); ok && p.tok == _Colon {
			return p.labeledStmtOrNil(label)
		}
		return p.simpleStmt(lhs, false)
	}

	switch p.tok {
	case _Lbrace:
		return p.blockStmt("")

	case _Var:
		return p.declStmt(p.varDecl)

	case _Operator, _Star:
		switch p.op {
		case Add, Sub, Mul, And, Xor, Not:
			return p.simpleStmt(nil, false) // unary operators
		}

	case _Literal, _Func, _Lparen, // operands
		_Lbrack, _Struct, _Map, _Chan, _Interface, // composite types
		_Arrow: // receive operator
		return p.simpleStmt(nil, false)

	case _For:
		return p.forStmt()

	case _If:
		return p.ifStmt()
	case _Parallel:
		return p.parallelStmt()

	case _Break, _Continue:
		s := new(BranchStmt)
		s.pos = p.pos()
		s.Tok = p.tok
		p.next()
		if p.tok == _Name {
			s.Label = p.name()
		}
		return s

	case _Go, _Defer:
		return p.callStmt()

	case _Goto:
		s := new(BranchStmt)
		s.pos = p.pos()
		s.Tok = _Goto
		p.next()
		s.Label = p.name()
		return s

	case _Return:
		s := new(ReturnStmt)
		s.pos = p.pos()
		p.next()
		if p.tok != _Semi && p.tok != _Rbrace {
			s.Results = p.exprList()
		}
		return s

	case _Semi:
		s := new(EmptyStmt)
		s.pos = p.pos()
		return s
	case _Wait:
		s := new(WaitStmt)
		s.pos = p.pos()
		p.next()
		return s
	}

	return nil
}

func (p *parser) labeledStmtOrNil(label *Name) Stmt {
	if trace {
		defer p.trace("labeledStmt")()
	}

	s := new(LabeledStmt)
	s.pos = p.pos()
	s.Label = label

	p.want(_Colon)

	if p.tok == _Rbrace {
		// We expect a statement (incl. an empty statement), which must be
		// terminated by a semicolon. Because semicolons may be omitted before
		// an _Rbrace, seeing an _Rbrace implies an empty statement.
		e := new(EmptyStmt)
		e.pos = p.pos()
		s.Stmt = e
		return s
	}

	s.Stmt = p.stmtOrNil()
	if s.Stmt != nil {
		return s
	}

	// report error at line of ':' token
	// p.syntaxErrorAt(s.pos, "missing statement after label")
	// we are already at the end of the labeled statement - no need to advance
	return nil // avoids follow-on errors (see e.g., fixedbugs/bug274.go)
}

// callStmt parses call-like statements that can be preceded by 'defer' and 'go'.
func (p *parser) callStmt() *CallStmt {
	if trace {
		defer p.trace("callStmt")()
	}

	s := new(CallStmt)
	s.pos = p.pos()
	s.Tok = p.tok // _Defer or _Go
	p.next()

	x := p.pexpr(p.tok == _Lparen) // keep_parens so we can report error below
	if t := unparen(x); t != x {
		p.error(fmt.Sprintf("expression in %s must not be parenthesized", s.Tok))
		// already progressed, no need to advance
		x = t
	}

	cx, ok := x.(*CallExpr)
	if !ok {
		p.error(fmt.Sprintf("expression in %s must be function call", s.Tok))
		// already progressed, no need to advance
		cx = new(CallExpr)
		cx.pos = x.Pos()
		cx.Fun = p.bad()
	}

	s.Call = cx
	return s
}

func (p *parser) header(keyword token) (init SimpleStmt, cond Expr, post SimpleStmt) {
	p.want(keyword)

	if p.tok == _Lbrace {
		if keyword == _If {
			// p.syntaxError("missing condition in if statement")
		}
		return
	}

	if func() bool { return true }() || true {

	}
	// p.tok != _Lbrace

	outer := p.xnest
	p.xnest = -1

	if p.tok != _Semi {
		// accept potential varDecl but complain
		if p.got(_Var) {
			// p.syntaxError(fmt.Sprintf("var declaration not allowed in %s initializer", keyword.String()))
		}
		init = p.simpleStmt(nil, keyword == _For)
		// If we have a range clause, we are done (can only happen for keyword == _For).
	}

	var condStmt SimpleStmt
	var semi struct {
		pos Pos
		lit string // valid if pos.IsKnown()
	}
	if p.tok != _Lbrace {
		if p.tok == _Semi {
			semi.pos = p.pos()
			semi.lit = p.lit
			p.next()
		} else {
			p.want(_Semi)
		}
		if keyword == _For {
			if p.tok != _Semi {
				if p.tok == _Lbrace {
					// p.syntaxError("expecting for loop condition")
					goto done
				}
				condStmt = p.simpleStmt(nil, false)
			}
			p.want(_Semi)
			if p.tok != _Lbrace {
				post = p.simpleStmt(nil, false)
				if a, _ := post.(*AssignStmt); a != nil && a.Op == Def {
					// p.syntaxError_at(a.Pos(), "cannot declare in post statement of for loop")
				}
			}
		} else if p.tok != _Lbrace {
			condStmt = p.simpleStmt(nil, false)
		}
	} else {
		condStmt = init
		init = nil
	}

done:
	// unpack condStmt
	switch s := condStmt.(type) {
	case nil:
		if keyword == _If /*&& semi.pos.IsKnown()*/ {
			if semi.lit != "semicolon" {
				// p.syntaxError_at(semi.pos, fmt.Sprintf("unexpected %s, expecting { after if clause", semi.lit))
			} else {
				// p.syntaxError_at(semi.pos, "missing condition in if statement")
			}
		}
	case *ExprStmt:
		cond = s.X
	default:
		// p.syntaxError(fmt.Sprintf("%s used as value", String(s)))
	}

	p.xnest = outer
	return
}

func (p *parser) ifStmt() *IfStmt {
	if trace {
		defer p.trace("ifStmt")()
	}

	s := new(IfStmt)
	s.pos = p.pos()
	s.Init, s.Cond, _ = p.header(_If)
	s.Then = p.blockStmt("if clause")
	// p.next()
	// fmt.Println(p.tok)
	// p.next()
	// fmt.Println(p.got(_Else))
	// fmt.Println(p.got(_Semi))
	// fmt.Println(p.got(_Else))
	// fmt.Println(p.tok)
	// fmt.Println(p.lit)

	if p.got(_Else) || (p.got(_Semi) && p.got(_Else)) {
		switch p.tok {
		case _If:
			s.Else = p.ifStmt()
		case _Lbrace:
			s.Else = p.blockStmt("")
		default:
			p.syntaxError("else must be followed by if or statement block")
			p.advance(_Name, _Rbrace)
		}
	}

	return s
}

func (p *parser) forStmt() Stmt {
	if trace {
		defer p.trace("forStmt")()
	}

	s := new(ForStmt)
	s.pos = p.pos()

	s.Init, s.Cond, s.Post = p.header(_For)
	s.Body = p.blockStmt("for clause")

	return s
}

func (p *parser) parallelStmt() Stmt {
	if trace {
		defer p.trace("parallelStmt")()
	}

	s := new(ParallelStmt)
	s.pos = p.pos()

	s.Body = p.blockStmt("for clause")
	return s
}

// SimpleStmt = EmptyStmt | ExpressionStmt | SendStmt | IncDecStmt | Assignment | ShortVarDecl .
func (p *parser) simpleStmt(lhs Expr, rangeOk bool) SimpleStmt {
	if trace {
		defer p.trace("simpleStmt")()
	}

	if lhs == nil {
		lhs = p.exprList()
	}

	if _, ok := lhs.(*ListExpr); !ok && p.tok != _Assign && p.tok != _Define {
		// expr
		pos := p.pos()
		switch p.tok {
		case _AssignOp:
			// lhs op= rhs
			op := p.op
			p.next()
			return p.newAssignStmt(pos, op, lhs, p.expr())

		case _IncOp:
			// lhs++ or lhs--
			op := p.op
			p.next()
			return p.newAssignStmt(pos, op, lhs, ImplicitOne)

		default:
			// expr
			s := new(ExprStmt)
			s.pos = lhs.Pos()
			s.X = lhs
			return s
		}
	}

	// expr_list
	pos := p.pos()
	switch p.tok {
	case _Assign:
		p.next()

		return p.newAssignStmt(pos, 0, lhs, p.exprList())

	case _Define:
		p.next()

		// expr_list ':=' expr_list
		rhs := p.exprList()
		as := p.newAssignStmt(pos, Def, lhs, rhs)
		return as

	default:
		p.syntaxError("expecting := or = or comma")
		p.advance(_Semi, _Rbrace)
		// make the best of what we have
		if x, ok := lhs.(*ListExpr); ok {
			lhs = x.ElemList[0]
		}
		s := new(ExprStmt)
		s.pos = lhs.Pos()
		s.X = lhs
		return s
	}
}

func (p *parser) newAssignStmt(pos Pos, op Operator, lhs, rhs Expr) *AssignStmt {
	a := new(AssignStmt)
	a.pos = pos
	a.Op = op
	a.Lhs = lhs
	a.Rhs = rhs
	return a
}

func (p *parser) funcBody() *BlockStmt {
	p.fnest++
	// errcnt := p.errcnt
	body := p.blockStmt("")
	p.fnest--

	// Don't check branches if there were syntax errors in the function
	// as it may lead to spurious errors (e.g., see test/switch2.go) or
	// possibly crashes due to incomplete syntax trees.
	// if p.mode&CheckBranches != 0 && errcnt == p.errcnt {
	checkBranches(body, p.errh)
	// }

	return body
}

func (p *parser) funcType() *FuncType {
	if trace {
		defer p.trace("funcType")()
	}

	typ := new(FuncType)
	typ.pos = p.pos()
	typ.ParamList = p.paramList()

	return typ
}

// IdentifierList = identifier { "," identifier } .
// The first name must be provided.
func (p *parser) nameList(first *Name) []*Name {
	if trace {
		defer p.trace("nameList")()
	}

	if first == nil {
		panic("first name not provided")
	}

	l := []*Name{first}
	for p.got(_Comma) {
		l = append(l, p.name())
	}

	return l
}

// usage: defer p.trace(msg)()
func (p *parser) trace(msg string) func() {
	p.print(msg + " (")
	const tab = ". "
	p.indent = append(p.indent, tab...)
	return func() {
		p.indent = p.indent[:len(p.indent)-len(tab)]
		if x := recover(); x != nil {
			panic(x) // skip print_trace
		}
		p.print(")")
	}
}

func (p *parser) print(msg string) {
	fmt.Printf("%5d: %s%s\n", p.line, p.indent, msg)
}

func (p *parser) got(tok token) bool {
	if p.tok == tok {
		p.next()
		return true
	}
	return false
}

func (p *parser) got2(tok token) bool {
	p.next()
	if p.tok == tok {
		p.next()
		return true
	}
	return false
}

// Convenience methods using the current token position.
// func (p *parser) pos() Pos         { return p.pos_at(p.line, p.col) }
// func (p *parser) error(msg string) { p.error_at(p.pos(), msg) }

func (p *parser) name() *Name {
	// no tracing to avoid overly verbose output

	if p.tok == _Name {
		n := p.newName(p.lit)
		p.next()
		return n
	}

	n := p.newName("_")
	p.syntaxError("expecting name")
	p.advance()
	return n
}

// appendGroup(f) = f | "(" { f ";" } ")" . // ";" is optional before ")"
func (p *parser) appendGroup(list []Decl, f func(*Group) Decl) []Decl {
	if p.tok == _Lparen {
		g := new(Group)
		p.list(_Lparen, _Semi, _Rparen, func() bool {
			list = append(list, f(g))
			return false
		})
	} else {
		list = append(list, f(nil))
	}

	return list
}

// ImportSpec = [ "." | PackageName ] ImportPath .
// ImportPath = string_lit .
func (p *parser) importDecl(group *Group) Decl {
	if trace {
		defer p.trace("importDecl")()
	}

	d := new(ImportDecl)
	d.pos = p.pos()
	if p.tok == _Name {
		d.LocalPkgName = p.name()
	}
	d.Path = p.oliteral()
	if d.Path == nil {
		p.syntaxError("missing import path")
		p.advance(_Semi, _Rparen)
		return nil
	}
	d.Group = group

	return d
}

func (p *parser) newName(value string) *Name {
	n := new(Name)
	n.pos = p.pos()
	n.Value = value
	return n
}

// Advance consumes tokens until it finds a token of the stopset or followlist.
// The stopset is only considered if we are inside a function (p.fnest > 0).
// The followlist is the list of valid tokens that can follow a production;
// if it is empty, exactly one (non-EOF) token is consumed to ensure progress.
func (p *parser) advance(followlist ...token) {
	if trace {
		p.print(fmt.Sprintf("advance %s", followlist))
	}

	// compute follow set
	// (not speed critical, advance is only called in error situations)
	var followset uint64 = 1 << _EOF // don't skip over EOF
	if len(followlist) > 0 {
		if p.fnest > 0 {
			followset |= stopset
		}
		for _, tok := range followlist {
			followset |= 1 << tok
		}
	}

	for !contains(followset, p.tok) {
		if trace {
			p.print("skip " + p.tok.String())
		}
		p.next()
		if len(followlist) == 0 {
			break
		}
	}

	if trace {
		p.print("next " + p.tok.String())
	}
}

func (p *parser) list(open, sep, close token, f func() bool) Pos {
	p.want(open)

	var done bool
	for p.tok != _EOF && p.tok != close && !done {
		done = f()
		// sep is optional before close
		if !p.got(sep) && p.tok != close {
			p.syntaxError(fmt.Sprintf("expecting %s or %s", tokstring(sep), tokstring(close)))
			p.advance(_Rparen, _Rbrack, _Rbrace)
			if p.tok != close {
				// position could be better but we had an error so we don't care
				return p.pos()
			}
		}
	}

	pos := p.pos()
	p.want(close)
	return pos
}

func (p *parser) pos() Pos               { return p.pos_at(p.line, p.col) }
func (p *parser) error(msg string)       { p.error_at(p.pos(), msg) }
func (p *parser) syntaxError(msg string) { p.syntaxErrorAt(p.pos(), msg) }

type lico uint32

// error reports an error at the given position.
func (p *parser) error_at(pos Pos, msg string) {
	// err := Error{pos, msg}
	// if p.first == nil {
	// 	p.first = err
	// }
	p.errcnt++
	if p.errh == nil {
		panic(p.first)
	}
	// p.errh(err)
}

// pos_at returns the Pos value for (line, col) and the current position base.
func (p *parser) pos_at(line, col uint) Pos {
	return MakePos(p.base, line, col)
}

// tokstring returns the English word for selected punctuation tokens
// for more readable error messages.
func tokstring(tok token) string {
	switch tok {
	case _Comma:
		return "comma"
	case _Semi:
		return "semicolon or newline"
	}
	return tok.String()
}

// unparen removes all parentheses around an expression.
func unparen(x Expr) Expr {
	for {
		p, ok := x.(*ParenExpr)
		if !ok {
			break
		}
		x = p.X
	}
	return x
}

// syntaxErrorAt reports a syntax error at the given position.
func (p *parser) syntaxErrorAt(pos Pos, msg string) {
	if trace {
		p.print("syntax error: " + msg)
	}

	if p.tok == _EOF && p.first != nil {
		return // avoid meaningless follow-up errors
	}

	// add punctuation etc. as needed to msg
	switch {
	case msg == "":
		// nothing to do
	case strings.HasPrefix(msg, "in "), strings.HasPrefix(msg, "at "), strings.HasPrefix(msg, "after "):
		msg = " " + msg
	case strings.HasPrefix(msg, "expecting "):
		msg = ", " + msg
	default:
		// plain error - we don't care about current token
		p.errorAt(pos, "syntax error: "+msg)
		return
	}

	// determine token string
	var tok string
	switch p.tok {
	case _Name, _Semi:
		tok = p.lit
	case _Literal:
		tok = "literal " + p.lit
	case _Operator:
		tok = p.op.String()
	case _AssignOp:
		tok = p.op.String() + "="
	case _IncOp:
		tok = p.op.String()
		tok += tok
	default:
		tok = tokstring(p.tok)
	}

	p.errorAt(pos, "syntax error: unexpected "+tok+msg)
}

// error reports an error at the given position.
func (p *parser) errorAt(pos Pos, msg string) {
	err := Error{pos, msg}
	if p.first == nil {
		p.first = err
	}
	p.errcnt++
	if p.errh == nil {
		panic(p.first)
	}
	p.errh(err)
}

// posAt returns the Pos value for (line, col) and the current position base.
func (p *parser) posAt(line, col uint) Pos {
	return MakePos(p.base, line, col)
}

func commentText(s string) string {
	if s[:2] == "/*" {
		return s[2 : len(s)-2] // lop off /* and */
	}

	// line comment (does not include newline)
	// (on Windows, the line comment may end in \r\n)
	i := len(s)
	if s[i-1] == '\r' {
		i--
	}
	return s[2:i] // lop off //, and \r at end, if any
}

// updateBase sets the current position base to a new line base at pos.
// The base's filename, line, and column values are extracted from text
// which is positioned at (tline, tcol) (only needed for error messages).
func (p *parser) updateBase(pos Pos, tline, tcol uint, text string) {
	i, n, ok := trailingDigits(text)
	if i == 0 {
		return // ignore (not a line directive)
	}
	// i > 0

	if !ok {
		// text has a suffix :xxx but xxx is not a number
		p.errorAt(p.posAt(tline, tcol+i), "invalid line number: "+text[i:])
		return
	}

	var line, col uint
	i2, n2, ok2 := trailingDigits(text[:i-1])
	if ok2 {
		//line filename:line:col
		i, i2 = i2, i
		line, col = n2, n
		if col == 0 || col > PosMax {
			p.errorAt(p.posAt(tline, tcol+i2), "invalid column number: "+text[i2:])
			return
		}
		text = text[:i2-1] // lop off ":col"
	} else {
		//line filename:line
		line = n
	}

	if line == 0 || line > PosMax {
		p.errorAt(p.posAt(tline, tcol+i), "invalid line number: "+text[i:])
		return
	}

	// If we have a column (//line filename:line:col form),
	// an empty filename means to use the previous filename.
	filename := text[:i-1] // lop off ":line"
	if filename == "" && ok2 {
		filename = p.base.Filename()
	}

	p.base = NewLineBase(pos, filename, line, col)
}

func trailingDigits(text string) (uint, uint, bool) {
	// Want to use LastIndexByte below but it's not defined in Go1.4 and bootstrap fails.
	i := strings.LastIndex(text, ":") // look from right (Windows filenames may contain ':')
	if i < 0 {
		return 0, 0, false // no ":"
	}
	// i >= 0
	n, err := strconv.ParseUint(text[i+1:], 10, 0)
	return uint(i + 1), uint(n), err == nil
}
