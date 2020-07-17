package tools

var IntToStringLevelError = map[int]string{
	PriorityNotify: "Notify",
	PriorityLow:    "Warning",
	PriorityMedium: "Error",
	PriorityHigh:   "Critical",
	PriorityInfo:   "Info",
}

func GetInfo(value int) string {
	if notify, ok := IntToStringLevelError[value]; ok {
		return notify
	}
	return "other"
}

const PriorityNotify = 0
const PriorityLow = 1
const PriorityMedium = 2
const PriorityHigh = 3
const PriorityInfo = 4
