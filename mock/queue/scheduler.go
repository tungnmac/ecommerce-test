package queue

import (
	"fmt"
	"log"
	"sync"
	"time"

	"go-mock/request"

	"github.com/robfig/cron/v3"
)

// Run scheduled job
func StartScheduler(apiURL string, batchSize int, totalRows int) {
	c := cron.New()

	// Run job every 1 minute
	c.AddFunc("@every 10s", func() {
		log.Println("Fetching data from API...")
		var wg sync.WaitGroup
		ch := make(chan []request.Data)

		// Start data fetching in parallel
		for i := 0; i < totalRows/batchSize; i++ {
			wg.Add(1)
			go func(page int) {
				defer wg.Done()
				request.FetchData(apiURL, batchSize, page, ch)
				time.Sleep(2 * time.Second)
			}(i)
		}

		// Read fetched data and push to queue
		go func() {
			for batch := range ch {
				PushToQueue(batch)
			}
		}()

		wg.Wait()
		close(ch)
		fmt.Println("Sync completed!")
	})

	c.Start()
	log.Println("Scheduler started...")
}
