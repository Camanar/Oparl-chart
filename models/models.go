package models

import "sync"

type Request struct {
	Url string
}

var (
	Queue   = make(chan Request, 100)
	Wg      sync.WaitGroup
	Workers = 5
)
