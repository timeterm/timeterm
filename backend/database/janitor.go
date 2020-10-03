package database

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	maxStopTime   = time.Second * 30
	maxJobRunTime = time.Minute
)

func (w *Wrapper) runJanitor(ctx context.Context) {
	c := cron.New(cron.WithLogger(w.logger))

	c.Schedule(&cron.ConstantDelaySchedule{
		Delay: time.Minute,
	}, cron.FuncJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), maxJobRunTime)
		defer cancel()

		err := w.DeleteOldOAuth2States(ctx)
		if err != nil {
			w.logger.Error(err, "could not delete old OAuth2 states")
		}
	}))

	c.Schedule(&cron.ConstantDelaySchedule{
		Delay: time.Minute,
	}, cron.FuncJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), maxJobRunTime)
		defer cancel()

		err := w.DeleteOldUserTokens(ctx)
		if err != nil {
			w.logger.Error(err, "could not delete old user tokens")
		}
	}))

	go c.Run()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), maxStopTime)
	defer cancel()

	select {
	case <-c.Stop().Done():
	case <-ctx.Done():
	}
}
