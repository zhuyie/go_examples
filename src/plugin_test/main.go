package main

import (
	"fmt"
	"plugin"
	"plugin_test/logic"
	"unsafe"
)

func main() {
	logicData := logic.GetLogicData("main")
	logicData.Value = 1234
	fmt.Printf("logicData addr=0x%x, Value=%v\n", unsafe.Pointer(logicData), logicData.Value)

	p, err := plugin.Open("./plugin/plugin.so")
	if err != nil {
		fmt.Printf("plugin.Open error=%v\n", err)
		return
	}
	testFunc, _ := p.Lookup("TestFunc")
	testFunc.(func())()

	fmt.Println("---------------------------------------------")

	logic.GetFuncAddr()

	testFunc2, _ := p.Lookup("TestFunc2")
	testFunc2.(func())()

	fmt.Println("---------------------------------------------")

	return
}
