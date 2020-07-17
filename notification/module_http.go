package notification

import (
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"time"
)

type ModuleHttp struct {
	Module
	LastTimeUpdate   time.Time
	Message          string
	ReadMessageChan  chan chan string
	BufferLog        []Message
	MaxBufferLog     int
	ReadMessagesChan chan chan []Message
}

func HttpTemplate(handlerFunc http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/favicon.ico" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		handlerFunc(w, r)
	})
}

func (ms *ModuleHttp) Init() (err error) {
	ms.MessageChan = make(chan Message, 1024)

	ms.ReadMessagesChan = make(chan chan []Message, 1024)
	ms.ReadMessageChan = make(chan chan string, 1024)

	ms.TimerChan = make(chan int)
	// Запускаю http сервер для ошибок
	//==============================================================
	// Получение ошибки
	http.Handle("/vitlog/error", HttpTemplate(func(w http.ResponseWriter, r *http.Request) {
		messChan := make(chan string)

		ms.ReadMessageChan <- messChan
		data := <-messChan
		//===========================================================
		w.WriteHeader(http.StatusOK)
		var err error
		if len(data) == 0 {
			_, err = fmt.Fprintf(w, "empty")
			return
		} else {
			_, err = fmt.Fprintf(w, data)
		}

		if err != nil {
			SendByService("default", fmt.Sprintf("%s - %s", ms.Name, err.Error()), tools.PriorityHigh)
		}
	}))
	// Получение ошибок
	http.Handle("/vitlog/log-json", HttpTemplate(func(w http.ResponseWriter, r *http.Request) {
		messChan := make(chan []Message)

		ms.ReadMessagesChan <- messChan
		data := <-messChan
		//===========================================================
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			SendByService("default", fmt.Sprintf("%s - %s", ms.Name, err.Error()), tools.PriorityHigh)
		}
	}))
	http.Handle("/vitlog/log", HttpTemplate(func(w http.ResponseWriter, r *http.Request) {
		messChan := make(chan []Message)

		ms.ReadMessagesChan <- messChan
		data := <-messChan

		//===========================================================
		w.WriteHeader(http.StatusOK)
		for _, item := range data {
			_, _ = fmt.Fprintf(w, "%s[%s] => %s\n", item.Date.Format("02.01.2006 15:04:05"), tools.GetInfo(item.Priority), item.Data)
		}
	}))

	//==============================================================
	go func() {
		if len(os.Getenv("VITLOG_HTTP_PORT")) == 0 {
			fmt.Printf("VitLog http server is not configured\n")
			return
		}
		host := os.Getenv("VITLOG_HTTP_HOST") + ":" + os.Getenv("VITLOG_HTTP_PORT")
		//==============================================================
		domain := os.Getenv("VITLOG_HTTP_HOST")
		if len(domain) == 0 {
			domain = "localhost"
		}
		fmt.Println("==========|http server check error|=========")
		println("http://" + domain + ":" + os.Getenv("VITLOG_HTTP_PORT"))
		server := &http.Server{Addr: host}

		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("VitLog : %s\n", err.Error())
		}
	}()
	go func() {
		for {
			select {
			case message := <-ms.MessageChan:
				ms.update(message)
			case message := <-ms.ReadMessageChan:
				message <- ms.Message
				// сбрасываем обшику через минут
				if time.Now().Sub(ms.LastTimeUpdate).Seconds() > 330 {
					ms.Message = ""
				}
			case message := <-ms.ReadMessagesChan:
				message <- ms.BufferLog
			}
		}
	}()
	return nil
}

func (ms *ModuleHttp) update(message Message) {
	ms.BufferLog = append(ms.BufferLog, message)

	length := len(ms.BufferLog)
	if length > ms.MaxBufferLog {
		ms.BufferLog = ms.BufferLog[length-ms.MaxBufferLog : length]
	}

	if tools.PriorityHigh > message.Priority || message.Priority == tools.PriorityInfo {
		return
	}

	ms.LastTimeUpdate = time.Now()
	if len(message.Data) > 125 {
		message.Data = message.Data[0:125]
	}

	ms.Message = fmt.Sprintf("%s err", message.Data)
	return
}

//===================================
func init() {
	AddModuleNotifications(&ModuleHttp{
		Module: Module{
			Name:     "http",
			Priority: tools.PriorityNotify,
		},
		MaxBufferLog: 100,
	})
}
