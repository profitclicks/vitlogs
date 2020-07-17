package vitlog

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error()
		return
	}

	err = Init(func(request *http.Request) {
		request.Header.Add("Authorization", "Basic "+os.Getenv("MASTER_TOKEN"))

	})

	if err != nil {
		fmt.Println(err.Error())
		t.Error()
		return
	}
	time.Sleep(time.Second)
	SendDebug("Hello world1")
	SendWarning("Hello world.0")
	SendError("Hello world3")
	SendCritical("Hello world4")
	SendInfo("Hello world5")

	host := os.Getenv("VITLOG_HTTP_HOST") + ":2122"
	//==============================================================
	domain := os.Getenv("VITLOG_HTTP_HOST")
	if len(domain) == 0 {
		domain = "localhost"
	}
	fmt.Println("==========|http server check error|=========")
	println("http://" + domain + ":2122")
	server := &http.Server{Addr: host}

	http.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendCritical("test error " + time.Now().String())

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "empty")
	}))
	err = server.ListenAndServe()

	fmt.Println(err)
	time.Sleep(time.Second * 20000000)

}
