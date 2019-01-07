package main

import (
	"fmt"
	"plugin"
	"plugin_globalvar/logic"
	"unsafe"
)

func main() {
	logicData := logic.GetLogicData()
	logicData.Value = 1234
	fmt.Printf("logicData addr=0x%x, Value=%v\n", unsafe.Pointer(logicData), logicData.Value)

	p, err := plugin.Open("./plugin.so")
	if err != nil {
		fmt.Printf("plugin.Open error=%v\n", err)
		return
	}
	testFunc, _ := p.Lookup("TestFunc")
	testFunc.(func())()

	return
}