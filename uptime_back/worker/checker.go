package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SohamChatterG/uptime/model"
	"github.com/SohamChatterG/uptime/repository"
	"github.com/SohamChatterG/uptime/service"
)

type Checker struct {
	urlRepo   *repository.URLRepository
	userRepo  *repository.UserRepository
	checkRepo *repository.CheckRepository
	notifySvc *service.GmailService
	interval  time.Duration
}

func NewChecker(urlRepo *repository.URLRepository, userRepo *repository.UserRepository, checkRepo *repository.CheckRepository, notifySvc *service.GmailService, interval time.Duration) *Checker { // <-- 2. CHANGE THIS from NotificationService
	return &Checker{
		urlRepo:   urlRepo,
		userRepo:  userRepo,
		checkRepo: checkRepo,
		notifySvc: notifySvc,
		interval:  interval,
	}
}

// Start, runChecks, and checkURL functions remain the same, but we'll add the alert logic.
func (c *Checker) Start() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	c.runChecks()
	for range ticker.C {
		c.runChecks()
	}
}

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

func (c *Checker) checkURL(url model.Url, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{Timeout: 10 * time.Second}
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
		check.WasSuccessful = false
		check.StatusCode = 0
	} else {
		defer resp.Body.Close()
		check.WasSuccessful = resp.StatusCode >= 200 && resp.StatusCode < 300
		check.StatusCode = resp.StatusCode
	}

	if err := c.checkRepo.Create(context.Background(), check); err != nil {
		log.Printf("Error saving check result for %s: %v", url.URL, err)
	}

	if url.Status != check.WasSuccessful {
		log.Printf("STATUS CHANGE: %s is now %s", url.URL, map[bool]string{true: "UP", false: "DOWN"}[check.WasSuccessful])

		user, err := c.userRepo.FindByID(context.Background(), url.UserID)
		if err != nil {
			log.Printf("Could not find user %s to send alert for URL %s", url.UserID.Hex(), url.Name)
		} else {
			var subject, message string
			if check.WasSuccessful {
				subject = fmt.Sprintf("✅ Resolved: Your site '%s' is back up!", url.Name)
				message = fmt.Sprintf("Good news! Your monitored URL '%s' (%s) has recovered and is now back online.", url.Name, url.URL)
			} else {
				subject = fmt.Sprintf("🔴 Alert: Your site '%s' is down!", url.Name)
				message = fmt.Sprintf("This is an automated alert to inform you that your monitored URL '%s' (%s) is currently down.", url.Name, url.URL)
			}
			// 3. Call the new SendNotification method
			c.notifySvc.SendNotification(user.Email, subject, message)
		}

		if err := c.urlRepo.UpdateStatus(context.Background(), url.ID, check.WasSuccessful); err != nil {
			log.Printf("Error updating URL status for %s: %v", url.URL, err)
		}
	}
}
