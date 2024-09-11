package stresstest

import "time"

type (
	Report struct {
		TimeSpent          time.Duration
		TotalRequests      int
		SuccessfulRequests int
		FailedRequests     map[int]int
	}
)
