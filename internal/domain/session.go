package domain

import "time"

type Session struct {
	Id      string
	Key     string
	Timeout time.Duration
}
