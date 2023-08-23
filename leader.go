package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/flowerinthenight/hedge"
	"github.com/golang/glog"
	gaxv2 "github.com/googleapis/gax-go/v2"
)

var (
	ctrlPingPong = "CTRL_PING_PONG"

	fnLeader = map[string]func(e *cloudevents.Event) ([]byte, error){
		ctrlPingPong: doLeaderPingPong,
	}
)

func leaderHandler(data interface{}, msg []byte) ([]byte, error) {
	var e cloudevents.Event
	err := json.Unmarshal(msg, &e)
	if err != nil {
		glog.Errorf("Unmarshal failed: %v", err)
		return nil, err
	}

	if _, ok := fnLeader[e.Type()]; !ok {
		return nil, fmt.Errorf("failed: unsupported type: %v", e.Type())
	}

	return fnLeader[e.Type()](&e)
}

func doLeaderPingPong(e *cloudevents.Event) ([]byte, error) {
	switch {
	case string(e.Data()) != "PING":
		return nil, fmt.Errorf("invalid message")
	default:
		return []byte("PONG"), nil
	}
}

func ensureLeaderActive(ctx context.Context) (bool, error) {
	msg := newEvent([]byte("PING"), "jupiter", ctrlPingPong)
	b, _ := json.Marshal(msg)
	r, err := sendToLeader(ctx, b)
	if err != nil {
		return false, err
	}

	switch {
	case string(r) == "PONG":
		return true, nil
	default:
		return false, nil
	}
}

func sendToLeader(ctx context.Context, m []byte) ([]byte, error) {
	result := make(chan []byte, 1)
	done := make(chan error, 1)
	go func() {
		var err error
		var res []byte
		defer func(b *[]byte, e *error) {
			result <- *b
			done <- *e
		}(&res, &err)

		bo := gaxv2.Backoff{Max: time.Minute}
		for i := 0; i < 10; i++ {
			if !op.IsRunning() {
				time.Sleep(bo.Pause())
				continue
			}
		}

		for i := 0; i < 10; i++ {
			var r []byte
			r, err = op.Send(ctx, m)
			if err != nil {
				time.Sleep(bo.Pause())
				continue
			}

			res = r // to outside
			return
		}
	}()

	for {
		select {
		case e := <-done:
			return <-result, e
		case <-ctx.Done():
			return nil, context.Canceled
		}
	}
}

func leaderLiveness(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	var active int32

	do := func() {
		atomic.StoreInt32(&active, 1)
		defer atomic.StoreInt32(&active, 0)
		hl, _ := op.HasLock()
		if !hl {
			return // leader's job only
		}

		b, _ := json.Marshal(newEvent(
			hedge.KeyValue{},
			"jupiter",
			ctrlBroadcastLeaderLiveness,
		))

		outs := op.Broadcast(ctx, b)
		for i, out := range outs {
			if out.Error != nil {
				glog.Errorf("leader liveness: broadcast[%v] failed: %v", i, out.Error)
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
		}

		if atomic.LoadInt32(&active) == 1 {
			continue
		}

		go do()
	}
}
