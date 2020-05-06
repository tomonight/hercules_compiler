package compile

import "hercules_compiler/engine/syntax"

const (
	bodyScope = ScopeID(1)
	funcScope = ScopeID(2)
)

var funcL map[string]*syntax.FuncDecl
var worldWideVar []*Var

func init() {
	funcL = make(map[string]*syntax.FuncDecl)
	worldWideVar = []*Var{}
}

//local func expr define

type Op uint8

const (
	ODELETE      Op = iota // delete(Left, Right)
	ODOT                   // Left.Sym (Left is of struct type)
	ODOTPTR                // Left.Sym (Left is of pointer to struct type)
	ODOTMETH               // Left.Sym (Left is non-interface, Right is method name)
	ODOTINTER              // Left.Sym (Left is interface, Right is method name)
	OXDOT                  // Left.Sym (before rewrite to one of the preceding)
	ODOTTYPE               // Left.Right or Left.Type (.Right during parsing, .Type once resolved); after walk, .Right contains address of interface type descriptor and .Right.Right contains address of concrete type descriptor
	ODOTTYPE2              // Left.Right or Left.Type (.Right during parsing, .Type once resolved; on rhs of OAS2DOTTYPE); after walk, .Right contains address of interface type descriptor
	OEQ                    // Left == Right
	ONE                    // Left != Right
	OLT                    // Left < Right
	OLE                    // Left <= Right
	OGE                    // Left >= Right
	OGT                    // Left > Right
	ODEREF                 // *Left
	OINDEX                 // Left[Right] (index of array or slice)
	OINDEXMAP              // Left[Right] (index of map)
	OKEY                   // Left:Right (key:value in struct/array/map literal)
	OSTRUCTKEY             // Sym:Left (key:value in struct literal, after type checking)
	OLEN                   // len(Left)
	OMUL                   // Left * Right
	ODIV                   // Left / Right
	OMOD                   // Left % Right
	OLSH                   // Left << Right
	ORSH                   // Left >> Right
	OAND                   // Left & Right
	OANDNOT                // Left &^ Right
	ONEW                   // new(Left)
	ONOT                   // !Left
	OBITNOT                // ^Left
	OPLUS                  // +Left
	ONEG                   // -Left
	OOROR                  // Left || Right
	OPANIC                 // panic(Left)
	OPRINT                 // print(List)
	OPRINTN                // println(List)
	OPAREN                 // (Left)
	OSEND                  // Left <- Right
	OSLICE                 // Left[List[0] : List[1]] (Left is untypechecked or slice)
	OSLICEARR              // Left[List[0] : List[1]] (Left is array)
	OSLICESTR              // Left[List[0] : List[1]] (Left is string)
	OSLICE3                // Left[List[0] : List[1] : List[2]] (Left is untypedchecked or slice)
	OSLICE3ARR             // Left[List[0] : List[1] : List[2]] (Left is array)
	OSLICEHEADER           // sliceheader{Left, List[0], List[1]} (Left is unsafe.Pointer, List[0] is length, List[1] is capacity)
	ORECOVER               // recover()
	ORECV                  // <-Left
	ORUNESTR               // Type(Left) (Type is string, Left is rune)
	OSELRECV               // Left = <-Right.Left: (appears as .Left of OCASE; Right.Op == ORECV)
	OSELRECV2              // List = <-Right.Left: (apperas as .Left of OCASE; count(List) == 2, Right.Op == ORECV)
	OIOTA                  // iota
	OALIGNOF               // unsafe.Alignof(Left)
	OOFFSETOF              // unsafe.Offsetof(Left)
	OSIZEOF                // unsafe.Sizeof(Left)
)

//define local function name
const (
	nprint      = "print"
	nlength     = "length"
	nappend     = "append"
	nsleep      = "sleep"
	nconnect    = "connect"
	nsetTarget  = "setTarget"
	nsetRunMode = "setRunmode"
)

type optype uint

const (
	opint    = optype(1)
	opstring = optype(2)
)

const (
	ostr uint = iota
	oint
	obool
	olist
	omap
)

//define engineer logic running status
const (
	EFF_CONTINUE int = iota
	EFF_STOP
	EFS_CONTINUE
	EFS_STOP

	//define callback result type
	CRT_PARSE_COMPELTED
	CRT_PARSE_FAILED
	// CRT_COMMAND_COMPLETED
	// CRT_COMMAND_FAILED
	// CRT_COMMAND_SKIPED
	CRT_STATEMENT_COMPLETED
	CRT_STATEMENT_FAILED
	CRT_STATEMENT_SKIPED
	CRT_SCRIPT_COMPLETED
	CRT_SCRIPT_FAILED
)
