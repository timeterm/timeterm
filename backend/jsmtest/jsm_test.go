package jsmtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gitlab.com/timeterm/timeterm/proto/go"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/proto"
)

func TestJSMGetMsg(t *testing.T) {
	conn, err := nats.Connect("localhost")
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

			bs, err := proto.Marshal(&timetermpb.DisownTokenMessage{
				DeviceId: uuid.New().String(),
			})
			if err != nil {
				t.Fatalf("Could not marshal DisownTokenMessage: %v", err)
			}

			_, err = conn.Request("FEDEV.ozuhLrexlBa4p50INjihAl.DISOWN-TOKEN", bs, time.Second*10)
			if err != nil {
				t.Fatalf("Could not publish message to stream: %v", err)
			}
		}()

		if i%1000 == 0 {
			fmt.Printf("%d/1000\n", i/1000)
		}
	}

	_ = sema.Acquire(context.Background(), 1000)
	fmt.Printf("1 mil messages produced in %s\n", time.Since(start))
}
