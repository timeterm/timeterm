package database

import (
	"context"
	"time"

	"github.com/go-logr/logr"
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
	}, newDeleteOldOAuth2StatesJob(w, w.logger))

	c.Schedule(&cron.ConstantDelaySchedule{
		Delay: time.Minute,
	}, newDeleteOldUserTokensJob(w, w.logger))

	c.Schedule(&cron.ConstantDelaySchedule{
		Delay: time.Minute,
	}, newDeleteOldDeviceTokensJob(w, w.logger))

	go c.Run()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), maxStopTime)
	defer cancel()

	select {
	case <-c.Stop().Done():
	case <-ctx.Done():
	}
}

type janitorJob struct {
	dbw    *Wrapper
	logger logr.Logger
	run    func(ctx context.Context, j janitorJob)
}

func (j janitorJob) Run() {
	if j.run != nil {
		ctx, cancel := context.WithTimeout(context.Background(), maxJobRunTime)
		defer cancel()

		j.run(ctx, j)
	}
}

func newJanitorJob(dbw *Wrapper, logger logr.Logger, run func(ctx context.Context, j janitorJob)) janitorJob {
	return janitorJob{
		dbw:    dbw,
		logger: logger,
		run:    run,
	}
}

func newDeleteOldUserTokensJob(dbw *Wrapper, logger logr.Logger) cron.Job {
	return newJanitorJob(dbw, logger, func(ctx context.Context, j janitorJob) {
		err := j.dbw.DeleteOldUserTokens(ctx)
		if err != nil {
			j.logger.Error(err, "could not delete old user tokens")
		}
	})
}

func newDeleteOldOAuth2StatesJob(dbw *Wrapper, logger logr.Logger) cron.Job {
	return newJanitorJob(dbw, logger, func(ctx context.Context, j janitorJob) {
		err := j.dbw.DeleteOldOAuth2States(ctx)
		if err != nil {
			j.logger.Error(err, "could not delete old OAuth2 states")
		}
	})
}

func newDeleteOldDeviceTokensJob(dbw *Wrapper, logger logr.Logger) cron.Job {
	return newJanitorJob(dbw, logger, func(ctx context.Context, j janitorJob) {
		err := j.dbw.DeleteOldDeviceTokens(ctx)
		if err != nil {
			j.logger.Error(err, "could not delete old device tokens")
		}
	})
}
