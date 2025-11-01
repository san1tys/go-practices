package main

import "time"

type Job struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Company   string    `json:"company"`
	Salary    int       `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
}

type JobsResponse struct {
	Items       []Job `json:"items"`
	NextAfterID *int  `json:"next_after_id,omitempty"`
}
