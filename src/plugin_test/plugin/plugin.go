package main

import (
	"fmt"
	"plugin_test/logic"
	"unsafe"
)

func TestFunc() {
	logicData := logic.GetLogicData("plugin")
	fmt.Printf("logicData addr=0x%x, Value=%v\n", unsafe.Pointer(logicData), logicData.Value)
}

func TestFunc2() {
	logic.GetFuncAddr()
}
