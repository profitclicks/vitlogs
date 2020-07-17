package notification

import (
	"fmt"
	"time"
)

var modules map[string]INotification

type Message struct {
	Data string
	Priority int
	Date time.Time
}

type Module struct {
	Name        string
	MessageChan chan Message
	TimerChan   chan int
	Priority    int
}
func (ms *Module) Send(message string, priority int) {
	if ms.Priority > priority {
		return
	}
	ms.MessageChan <- Message{
		Data:     message,
		Priority: priority,
		Date: time.Now(),
	}
}

func (ms *Module) GetName() string {
	return ms.Name
}

type INotification interface {
	Init() error
	Send(message string, priority int)
	GetName() string
}

func Notifications(body string, priority int) {
	for _, module := range modules {
		module.Send(body, priority)
	}
}
func SendByService(name string,body string, priority int) {
	if module, ok := modules[name]; ok{
		module.Send(body, priority)
	}
}


func InitNotifications() error {
	fmt.Println("==========|notification|=========")
	for name, module := range modules {
		fmt.Printf("notification module %s init\n", name)
		if err := module.Init(); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}
// TODO хрень но лучше не придумал
func AddModuleNotifications(module INotification) {
	if modules == nil {
		modules = make( map[string]INotification)
	}
	modules[module.GetName()] = module
}
func RemoveModuleNotifications(name string) {
	fmt.Printf("notification module %s disable\n", name)
	delete(modules, name)
}