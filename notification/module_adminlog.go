package notification

import (
	"fmt"
	"github.com/profitclicks/vitlogs/tools"

	"strings"
	"time"
)

type ModuleAdminLog struct {
	Module
	Items []string
}

func (ms *ModuleAdminLog) Init() (err error) {
	ms.MessageChan = make(chan Message, 1024)
	ms.TimerChan = make(chan int)

	go func() {
		for {
			select {
			case message := <-ms.MessageChan:
				ms.add(message)
			case <-ms.TimerChan:
				ms.send()
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			ms.TimerChan <- 0
		}
	}()

	return nil
}
func (ms *ModuleAdminLog) add(message Message) {
	date := time.Now()
	ms.Items = append(ms.Items, fmt.Sprintf("%s <%s> %s<eof>", date.Format("2006.01.02 15:04:05"), tools.GetInfo(message.Priority), message.Data))
	return
}
func (ms *ModuleAdminLog) send() {
	if len(ms.Items) > 0 {
		fmt.Println(strings.Join(ms.Items, "\n"))
		ms.Items = nil
	}
	return
}

//===================================
func init() {
	AddModuleNotifications(&ModuleAdminLog{
		Module: Module{
			Name:     "default",
			Priority: tools.PriorityNotify,
		},
	})
}
