package logic

import (
	"fmt"
	"reflect"
	"unsafe"
)

type LogicData struct {
	Value int64
}

var theLogicData LogicData

func GetLogicData(caller string) *LogicData {
	fmt.Printf("GetLogicData caller=%v\n", caller)
	return &theLogicData
}

type funcval struct {
	fn uintptr
	// variable-size, fn-specific data here
}

type value struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func foo(v bool) {
	if v {
		fmt.Println("foo true")
	} else {
		fmt.Println("foo false")
	}
}

func GetFuncAddr() {
	foo(true)

	f := foo
	f(false)
	// f 实质是一个 *funcval
	fv := *(**funcval)(unsafe.Pointer(&f))
	fmt.Printf("func_addr = 0x%x\n", fv.fn)

	v := reflect.ValueOf(foo)
	v.Call([]reflect.Value{reflect.ValueOf(true)})
	// v 实质是一个 value
	vv := (*value)(unsafe.Pointer(&v))
	// 对于Func类型的value，其ptr指向一个*funcval
	fv2 := (*funcval)(vv.ptr)
	fmt.Printf("func_addr = 0x%x\n", fv2.fn)
}
