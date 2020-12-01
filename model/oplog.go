package model

import "time"

type OpLog struct {
	Ts       time.Time
	OP       string
	Object   string
	Content  string
	Duration time.Duration
}
