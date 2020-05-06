package compile

import "fmt"

type T uint

const (
	intT T = iota
	stringT
	listT
	mapT
)

func checkTarget(params map[string]interface{}) (*Target, error) {
	target := &Target{}
	if assertType(params["protocol"], stringT) {
		target.protocol = params["protocol"].(string)
	}
	if assertType(params["host"], stringT) {
		target.host = params["host"].(string)
	}
	if assertType(params["username"], stringT) {
		target.username = params["username"].(string)
	}
	if assertType(params["password"], stringT) {
		target.password = params["password"].(string)
	}
	if assertType(params["port"], intT) {
		target.port = params["port"].(int)
	}
	if assertType(params["keyfile"], stringT) {
		target.keyfile = params["keyfile"].(string)
	}

	return target, nil
}

// func targetValid(target *Target) error {
// 	if target.keyfile != "" {
// 		return nil
// 	}
// 	if target.host == "" || target.password == "" || target.port == 0 || target.protocol == "" || target.username == "" {
// 		return fmt.Errorf("target params invalid")
// 	}
// 	return nil
// }

func assertType(k interface{}, t T) bool {
	if k == nil {
		return false
	}
	ok := false
	switch t {
	case intT:
		_, ok = k.(int)
	case stringT:
		_, ok = k.(string)
	default:
		return false
	}
	return ok
}

func getBasicString(k interface{}) string {
	return fmt.Sprintf("%v", k)
}
