package notification

import (
	"bytes"
	"fmt"
	"github.com/profitclicks/vitlogs/tools"

	"net/http"
	"os"
	"strings"
	"time"
)

type ModuleSlack struct {
	Module
	Items []string
}

func (ms *ModuleSlack) Init() (err error) {
	if len(os.Getenv("SLACK_KEY")) == 0 {
		RemoveModuleNotifications(ms.Name)
		return
	}

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
	return
}
func (ms *ModuleSlack) add(message Message) {
	ms.Items = append(ms.Items, fmt.Sprintf("%s <%s> %s", message.Date.Format("2006.01.02 15:04:05"), tools.GetInfo(message.Priority), message.Data))
	return
}

func (ms *ModuleSlack) send() {
	if len(ms.Items) == 0 {
		return
	}
	message := strings.Join(ms.Items, "\n")
	isMultiLine := false

	if len(message) > 1 {
		isMultiLine = true
	}
	ms.Items = nil

	format := "%s : %s"
	if isMultiLine {
		format = "Service '%s' : \n%s"
	}

	data := fmt.Sprintf(format, os.Getenv("NAME"), message)

	url := fmt.Sprintf("https://hooks.slack.com/services/%s", os.Getenv("SLACK_KEY"))

	var jsonStr = []byte(fmt.Sprintf(`{"text":"%s"}`, data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		SendByService("default", fmt.Sprintf("%s - %s", ms.Name, err.Error()), tools.PriorityHigh)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		SendByService("default", fmt.Sprintf("%s - %s", ms.Name, err.Error()), tools.PriorityHigh)
		return
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			SendByService("default", fmt.Sprintf("%s - %s", ms.Name, err.Error()), tools.PriorityHigh)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		SendByService("default", fmt.Sprintf("%s - %s", ms.Name, "http status is not valid"), tools.PriorityHigh)
	}
}

//===================================
func init() {
	AddModuleNotifications(&ModuleSlack{
		Module: Module{
			Name:     "slack",
			Priority: tools.PriorityMedium,
		},
	})
}
