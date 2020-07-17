package vitlog

import (
	"github.com/profitclicks/vitlogs/notification"
	"github.com/profitclicks/vitlogs/tools"
	"net/http"
	"os"
	"strings"
	"time"
)

var isDebug bool

func Init(request func(request *http.Request)) error {
	if os.Getenv("DEBUG") == "1" {
		isDebug = true
	}
	if err := notification.InitNotifications(); err != nil {
		return err
	}

	// Оповещаем масте о присутствии
	go func() {
		if len(os.Getenv("GROUP")) == 0 || len(os.Getenv("MASTER_DOMAIN")) == 0 {
			SendNotify("!!! presence alert disabled !!!")
			return
		}

		for {
			req, err := http.NewRequest("GET", os.Getenv(`TRACK_URL`), nil)

			if err != nil {
				SendWarning("vitchecker", err.Error())
				continue
			}

			if request != nil {
				request(req)
			}

			req.Header.Add("Accept", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				SendWarning("vitchecker", err.Error())
				time.Sleep(time.Second * 5)
				continue
			}
			if err := resp.Body.Close(); err != nil {
				SendWarning("vitchecker", err.Error())
				time.Sleep(time.Second * 5)
			}
			time.Sleep(time.Minute)
		}
	}()
	return nil
}
func SendDebug(body ...string) {
	if !isDebug {
		return
	}
	notification.Notifications(strings.Join(body, " "), tools.PriorityNotify)
}
func SendNotify(body ...string) {
	notification.Notifications(strings.Join(body, " "), tools.PriorityNotify)
}
func SendWarning(body ...string) {
	notification.Notifications(strings.Join(body, " "), tools.PriorityLow)
}
func SendError(body ...string) {
	notification.Notifications(strings.Join(body, " "), tools.PriorityMedium)
}
func SendCritical(body ...string) {
	notification.Notifications(strings.Join(body, " "), tools.PriorityHigh)
}
func SendInfo(body ...string) {
	notification.Notifications(strings.Join(body, " "), tools.PriorityInfo)
}
