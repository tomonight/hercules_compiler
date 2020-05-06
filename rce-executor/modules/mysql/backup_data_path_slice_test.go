package mysql

import (
	"sort"
	"testing"
)

func TestBackUpDataPathSlice(t *testing.T) {
	pathList := BackUpDataPathSlice{"20190121163200", "1", "20190122144300"}
	t.Log("pathlist valid = ", pathList.Valid())
	sort.Sort(pathList)
	t.Log("pathlist=", pathList)
}
