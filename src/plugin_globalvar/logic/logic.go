package logic

type LogicData struct {
	Value int64
}

var theLogicData LogicData

func GetLogicData() *LogicData {
	return &theLogicData
}
