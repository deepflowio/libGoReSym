package main

import "C"
import (
	"bytes"
	"debug/buildinfo"
	"io/ioutil"
	"strings"

	"fmt"

	"github.com/mandiant/GoReSym/objfile"
)

func main() {

}

//export function_address
func function_address(file_name *C.char, func_name *C.char) (addr uintptr, size int) {
	fileName := C.GoString(file_name)
	funcName := C.GoString(func_name)

	versionOverride := ""

	file, err := objfile.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	tab, _, err := file.PCLineTable(versionOverride)
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

//export itab_address
func itab_address(file_name *C.char, itab_name *C.char) (addr uintptr) {
	fileName := C.GoString(file_name)
	iTabName := C.GoString(itab_name)

	runtimeVersion := ""

	file, err := objfile.Open(fileName)
	if err != nil {
		return
	}

	// try to get version the 'correct' way, also fill out buildSettings if parsing was ok
	bi, err := buildinfo.ReadFile(fileName)
	if err == nil {
		runtimeVersion = bi.GoVersion
	}

	fileData, fileDataErr := ioutil.ReadFile(fileName)
	if fileDataErr != nil {
		return
	}
	// GOVERSION
	if runtimeVersion == "" {
		// go1.<varies><garbage data>
		idx := bytes.Index(fileData, []byte{0x67, 0x6F, 0x31, 0x2E})
		if idx != -1 && len(fileData[idx:]) > 10 {
			runtimeVersion = "go1."
			ver := fileData[idx+4 : idx+10]
			for i, c := range ver {
				// the string is _not_ null terminated, nor length delimited. So, filter till first non-numeric ascii
				nextIsNumeric := (i+1) < len(ver) && ver[i+1] >= 0x30 && ver[i+1] <= 0x39

				// careful not to end with a . at the end
				if (c >= 0x30 && c <= 0x39 && c != ' ') || (c == '.' && nextIsNumeric) {
					runtimeVersion += string([]byte{c})
				} else {
					break
				}
			}
		}
	}

	tab, tabva, err := file.PCLineTable(runtimeVersion)
	if err != nil {
		return
	}

	if tab.Go12line == nil {
		return
	}

	// numeric only, go1.17 -> 1.17
	goVersionIdx := strings.Index(runtimeVersion, "go")
	if goVersionIdx != -1 {
		// "devel go1.18-2d1d548 Tue Dec 21 03:55:43 2021 +0000"
		runtimeVersion = strings.Split(runtimeVersion[goVersionIdx+2:]+" ", " ")[0]

		// go1.18-2d1d548
		runtimeVersion = strings.Split(runtimeVersion+"-", "-")[0]
	}

	is64bit := tab.Go12line.Ptrsize == 8
	littleendian := tab.Go12line.Binary.String() == "LittleEndian"

	// this can be a little tricky to locate and parse properly across all go versions
	_, moduleData, err := file.ModuleDataTable(tabva, runtimeVersion, tab.Go12line.Version.String(), is64bit, littleendian)
	if err != nil {
		return
	}
	entry := file.Entries()[0]

	parts := strings.Split(runtimeVersion, ".")
	if len(parts) >= 2 {
		runtimeVersion = parts[0] + "." + parts[1]
	}

	ptrSize := uint64(0)
	if is64bit {
		ptrSize = 8
	} else {
		ptrSize = 4
	}

	for i := 0; i < int(moduleData.ITablinks.Len); i++ {
		itabAddr, err := entry.ReadPointerSizeMem(uint64(moduleData.ITablinks.Data)+ptrSize*uint64(i), is64bit, littleendian)
		if err != nil {
			continue
		}

		interfaceAddr, err := entry.ReadPointerSizeMem(itabAddr, is64bit, littleendian)
		if err != nil {
			continue
		}

		typeAddr, err := entry.ReadPointerSizeMem(itabAddr+ptrSize, is64bit, littleendian)
		if err != nil {
			continue
		}

		parsed, err := entry.ParseType(runtimeVersion, moduleData, interfaceAddr, is64bit, littleendian)
		parsed2, err2 := entry.ParseType(runtimeVersion, moduleData, typeAddr, is64bit, littleendian)
		if err == nil && err2 == nil && len(parsed) > 0 && len(parsed2) > 0 {
			interfaceName := parsed[0].Str
			implementerName := parsed2[0].Str
			if fmt.Sprintf("go.itab.%s,%s", implementerName, interfaceName) != iTabName {
				continue
			}
			return uintptr(itabAddr)
		}
	}
	return
}
