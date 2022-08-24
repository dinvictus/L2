package main

import (
	"testing"
	"time"
)

func TestChannelDone(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(1*time.Second),
		sig(2000*time.Millisecond),
		sig(5000*time.Millisecond),
		sig(3*time.Second),
		sig(4*time.Second),
	)
	if int(time.Since(start).Seconds()-5) != 0 {
		t.Fatal("Error done channel")
	}
}
