package worker

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SohamChatterG/uptime/model"
	"github.com/SohamChatterG/uptime/repository"
)

type Checker struct {
	urlRepo   *repository.URLRepository
	checkRepo *repository.CheckRepository
	interval  time.Duration
}

func NewChecker(urlRepo *repository.URLRepository, checkRepo *repository.CheckRepository, interval time.Duration) *Checker {
	return &Checker{
		urlRepo:   urlRepo,
		checkRepo: checkRepo,
		interval:  interval,
	}
}

// Start begins the periodic checking process.
func (c *Checker) Start() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// Run once immediately on start, then wait for the ticker.
	c.runChecks()

	for range ticker.C {
		c.runChecks()
	}
}

// runChecks fetches all active URLs and checks them concurrently.
func (c *Checker) runChecks() {
	log.Println("Running uptime checks...")
	urls, err := c.urlRepo.GetAllActive(context.Background())
	if err != nil {
		log.Printf("Error fetching URLs for checking: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go c.checkURL(url, &wg)
	}
	wg.Wait()
	log.Println("Uptime checks finished.")
}

// checkURL performs an HTTP request for a single URL and logs the result.
func (c *Checker) checkURL(url model.Url, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url.URL)
	duration := time.Since(start).Milliseconds()

	check := &model.Check{
		UrlID:          url.ID,
		UserID:         url.UserID,
		CheckedAt:      time.Now(),
		ResponseTimeMS: duration,
	}

	if err != nil {
		log.Printf("URL DOWN: %s (Error: %v)", url.URL, err)
		check.WasSuccessful = false
		check.StatusCode = 0
	} else {
		defer resp.Body.Close()
		log.Printf("URL CHECK: %s (Status: %s)", url.URL, resp.Status)
		check.WasSuccessful = resp.StatusCode >= 200 && resp.StatusCode < 300
		check.StatusCode = resp.StatusCode
	}

	// Create the check record in the database
	if err := c.checkRepo.Create(context.Background(), check); err != nil {
		log.Printf("Error saving check result for %s: %v", url.URL, err)
	}

	// Update the URL's last known status if it has changed
	if url.Status != check.WasSuccessful {
		log.Printf("STATUS CHANGE: %s is now %s", url.URL, map[bool]string{true: "UP", false: "DOWN"}[check.WasSuccessful])
		if err := c.urlRepo.UpdateStatus(context.Background(), url.ID, check.WasSuccessful); err != nil {
			log.Printf("Error updating URL status for %s: %v", url.URL, err)
		}
		// In a real application, you would trigger an email/webhook alert here!
	}
}
