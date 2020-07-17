package notification

import (
	"fmt"

	"log/syslog"
	"os"
)

type ModuleSyslog struct {
	Module
	//Syslog *syslog.Client
	Syslog *syslog.Writer
}

func (ms *ModuleSyslog) Init() (err error) {
	if len(os.Getenv("SYSLOG_HOST")) == 0 {
		RemoveModuleNotifications(ms.Name)
		return
	}

	ms.Syslog, err = syslog.Dial("udp", os.Getenv("SYSLOG_HOST"),
		syslog.LOG_WARNING|syslog.LOG_LOCAL0, os.Getenv("NAME"))
	if err != nil {
		return
	}
	ms.MessageChan = make(chan Message, 1024)

	go func() {
		for {
			select {
			case message := <-ms.MessageChan:
				ms.send(message)
			}
		}
	}()

	return nil
}
func (ms *ModuleSyslog) send(message Message) {
	data := fmt.Sprintf("<%s> %s", tools.IntToStringLevelError[message.Priority], message.Data)
	switch message.Priority {
	case tools.PriorityLow:
		_ = ms.Syslog.Warning(data)
	case tools.PriorityMedium:
		_ = ms.Syslog.Err(data)
	case tools.PriorityHigh:
		_ = ms.Syslog.Crit(data)
	default:
		_ = ms.Syslog.Notice(data)
	}
	return
}

//===================================
func init() {
	AddModuleNotifications(&ModuleSyslog{
		Module: Module{
			Name:     "syslog",
			Priority: tools.PriorityNotify,
		},
	})
}
