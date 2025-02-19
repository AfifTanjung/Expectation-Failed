package scheduler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

func hitApi() {
	resp, err := http.Post("https://telkom-service-319971276127.asia-southeast2.run.app/api/v1/validation/process-batch", "application/json", nil)
	if err != nil {
		fmt.Println("Error Hit API:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("API Hit successfully at", time.Now(), "Status:", resp.StatusCode)
}

func StartScheduler() *cron.Cron {
	c := cron.New()
	c.AddFunc("@every 1m", hitApi)
	c.Start()
	fmt.Println("Scheduler Started...")
	return c
}
