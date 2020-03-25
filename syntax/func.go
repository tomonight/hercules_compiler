package syntax

//FuncList load all function is memory
//symboltable
var FuncList map[string]*FuncDecl

func init() {
	FuncList = make(map[string]*FuncDecl)
}

func addFunc(f *FuncDecl) error {

	if _, ok := FuncList[f.Name.Value]; ok {
		return functionRedeclireError
	}
	FuncList[f.Name.Value] = f
	return nil
}

//GetFunc getFunc Declire by func name
func GetFunc(name string) *FuncDecl {
	return FuncList[name]
}

//FuncExsit check if name of the function is exsit
func FuncExsit(name string) bool {
	if _, ok := FuncList[name]; ok {
		return true
	}
	return false
}
