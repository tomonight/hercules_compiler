package syntax

//FuncList load all function is memory
//symboltable
var FuncList map[string]*FuncDecl

var ConstList []*ConstDecl

func init() {
	FuncList = make(map[string]*FuncDecl)
	ConstList = []*ConstDecl{}
}

func addFunc(f *FuncDecl) error {

	if _, ok := FuncList[f.Name.Value]; ok {
		return functionRedeclireError
	}
	FuncList[f.Name.Value] = f
	return nil
}

func addConstDecl(decl *ConstDecl) error {
	ConstList = append(ConstList, decl)
	return nil
}

func CheckIsConst(name string) bool {
	for _, v := range ConstList {
		if v.NameList[0].Value == name {
			return true
		}
	}
	return false
}

func GetConst(name string) *ConstDecl {
	for _, v := range ConstList {
		if v.NameList[0].Value == name {
			return v
		}
	}
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
