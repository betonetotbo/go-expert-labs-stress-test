package stresstest

import (
	"log"
	"math"
	"time"
)

func StressTest(url, method string, headers []string, body []byte, timeout time.Duration, requests, concurrency int) *Report {
	requestsPerBatch := int(math.Floor(float64(requests) / float64(concurrency)))

	start := time.Now()
	reports := make(chan Report)
	for i := 0; i < concurrency; i++ {
		if i == concurrency-1 && requestsPerBatch*concurrency < requests {
			requestsPerBatch++
		}

		log.Printf("Requests per batch: %d, concurrency: %d", requestsPerBatch, concurrency)

		go func(count int) {
			r := Report{
				FailedRequests: make(map[int]int),
			}
			for i := 0; i < count; i++ {
				code := makeRequest(url, method, headers, body, timeout)
				if code >= 200 && code < 300 {
					r.SuccessfulRequests++
				} else {
					r.FailedRequests[code]++
				}
			}
			reports <- r
		}(requestsPerBatch)
	}

	report := Report{
		FailedRequests: make(map[int]int),
	}
	for i := 0; i < concurrency; i++ {
		r := <-reports
		report.SuccessfulRequests += r.SuccessfulRequests
		report.TotalRequests += r.SuccessfulRequests
		for k, v := range r.FailedRequests {
			report.FailedRequests[k] += v
			if k > 0 {
				report.TotalRequests += v
			}
		}
	}

	report.TimeSpent = time.Since(start)
	return &report
}
