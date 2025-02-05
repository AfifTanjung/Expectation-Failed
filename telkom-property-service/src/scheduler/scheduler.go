package scheduler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

func hitApi() {
	resp, err := http.Post("http://localhost:3001/validations/validate", "application/json", nil)
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
