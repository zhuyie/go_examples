package logic

import "fmt"

type LogicData struct {
	Value int64
}

var theLogicData LogicData

func GetLogicData(caller string) *LogicData {
	fmt.Printf("GetLogicData caller=%v\n", caller)
	return &theLogicData
}
