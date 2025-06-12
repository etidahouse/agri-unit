package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Name     string          `json:"name"`
	Schedule string          `json:"schedule"`
	Method   string          `json:"method"`
	URL      string          `json:"url"`
	Body     json.RawMessage `json:"body,omitempty"`
}

func runJob(job Job) {
	log.Printf("Running job: %s", job.Name)

	var body io.Reader
	if len(job.Body) > 0 {
		body = bytes.NewBuffer(job.Body)
	}

	req, err := http.NewRequest(job.Method, job.URL, body)
	if err != nil {
		log.Printf("[%s] Error creating request: %v", job.Name, err)
		return
	}

	if job.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[%s] HTTP error: %v", job.Name, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[%s] Status: %s", job.Name, resp.Status)
}

func main() {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "/config/jobs.json"
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var jobs []Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	c := cron.New()

	for _, job := range jobs {
		jobCopy := job
		_, err := c.AddFunc(job.Schedule, func() {
			runJob(jobCopy)
		})
		if err != nil {
			log.Printf("Error adding job %s: %v", job.Name, err)
		} else {
			log.Printf("Scheduled job %s", job.Name)
		}
	}

	c.Run()
}
