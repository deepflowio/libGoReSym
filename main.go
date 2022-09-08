package main

import (
	"C"

	"github.com/mandiant/GoReSym/objfile"
)

func main() {

}

//export FunctionAddress
func FunctionAddress(fileName string, funcName string) (addr uintptr, size int) {
	file, err := objfile.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	tab, _, err := file.PCLineTable("")
	if err != nil {
		return
	}

	if tab.Go12line == nil {
		return
	}

	for _, elem := range tab.Funcs {
		if elem.Name != funcName {
			continue
		}
		addr = uintptr(elem.Entry)
		size = int(elem.End - elem.Entry)
		break
	}
	return
}
