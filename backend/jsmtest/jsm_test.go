package jsmtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"golang.org/x/sync/semaphore"
)

func TestJSMGetMsg(t *testing.T) {
	conn, err := nats.Connect("localhost", nats.UseOldRequestStyle())
	if err != nil {
		t.Fatalf("Could not connect to NATS: %v", err)
	}

	sema := semaphore.NewWeighted(50000)
	start := time.Now()

	for i := 0; i < 1000000; i++ {
		err := sema.Acquire(context.Background(), 1)
		if err != nil {
			t.Fatalf("Could not acquire semaphore with weight 1: %v", err)
		}

		go func() {
			defer sema.Release(1)

			_, err = conn.Request("FEDEV.ozuhLrexlBa4p50INjihAl.DISOWN-TOKEN", []byte("TEST"), time.Second*10)
			if err != nil {
				t.Fatalf("Could not publish message to stream: %v", err)
			}
		}()

		if i % 1000 == 0 {
			fmt.Printf("%d/1000\n", i / 1000)
		}
	}

	sema.Acquire(context.Background(), 1000)
	fmt.Printf("1 mil messages produced in %s\n", time.Since(start))
}
