package main

import (
	"fmt"
	"net/http"
	"time"
)

// Site represents a website to check
type Site struct {
	URL string
}

// Result holds the status of a checked site
type Result struct {
	URL     string
	Status  string
	Latency time.Duration
	Error   error
}

func main() {
	// A list of websites to check concurrently
	links := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://go.dev",
		"https://this-is-a-fake-site.com",
	}

	// Create a channel to communicate the Results
	// Channels allow goroutines to send data to each other safely
	resultsChan := make(chan Result)

	fmt.Println("Starting site health check...")
	fmt.Println("-------------------------------")

	// Launch a goroutine for each link
	for _, link := range links {
		go checkLink(link, resultsChan)
	}

	// Receive the results from the channel
	// Since we know we launched len(links) goroutines, we loop that many times
	for i := 0; i < len(links); i++ {
		res := <-resultsChan // This blocks until a result is sent into the channel
		displayResult(res)
	}

	fmt.Println("-------------------------------")
	fmt.Println("Check complete.")
}

// checkLink performs an HTTP GET request and sends the result to a channel
func checkLink(link string, c chan Result) {
	start := time.Now()
	
	resp, err := http.Get(link)
	elapsed := time.Since(start)

	if err != nil {
		c <- Result{URL: link, Error: err, Latency: elapsed}
		return
	}
	defer resp.Body.Close()

	c <- Result{
		URL:     link,
		Status:  resp.Status,
		Latency: elapsed,
		Error:   nil,
	}
}

// displayResult prints the outcome to the console
func displayResult(r Result) {
	if r.Error != nil {
		fmt.Printf("[DOWN] %-30s | Error: %v\n", r.URL, r.Error)
	} else {
		fmt.Printf("[UP]   %-30s | Status: %-15s | Took: %v\n", r.URL, r.Status, r.Latency)
	}
}
