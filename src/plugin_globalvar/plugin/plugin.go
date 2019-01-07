package main

import (
	"fmt"
	"plugin_globalvar/logic"
	"unsafe"
)

func TestFunc() {
	logicData := logic.GetLogicData()
	fmt.Printf("logicData addr=0x%x, Value=%v\n", unsafe.Pointer(logicData), logicData.Value)
}
