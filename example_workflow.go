package example

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// ExampleWorkflow workflow definition
func ExampleWorkflow(ctx workflow.Context) error {

	logger := workflow.GetLogger(ctx)

	stop := false

	initChannel := workflow.GetSignalChannel(ctx, "init")
	resetChannel := workflow.GetSignalChannel(ctx, "stop")

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(initChannel, func(c workflow.ReceiveChannel, more bool) {
			var req interface{}
			c.Receive(ctx, &req)

			ao := workflow.ActivityOptions{
				ScheduleToCloseTimeout: 24 * time.Hour,
				StartToCloseTimeout:    24 * time.Hour,

				HeartbeatTimeout: 10 * time.Second,
			}

			ctx = workflow.WithActivityOptions(ctx, ao)
			res := workflow.ExecuteActivity(ctx, ExampleActivity)
			selector.AddFuture(res, func(f workflow.Future) {
				err := f.Get(ctx, nil)
				if err != nil {
					logger.Error(err.Error())
					return
				}
			})

			logger.Info("done with init")
		})

		selector.AddReceive(resetChannel, func(c workflow.ReceiveChannel, more bool) {
			var req interface{}
			c.Receive(ctx, &req)

			stop = true
		})

		selector.Select(ctx)

		if stop {
			break
		}
	}

	workflow.GetLogger(ctx).Info("Workflow completed.")
	return nil
}

func ExampleActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)

	logger.Info("invoking activity")

	t := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-t.C:
				logger.Info("heartbeating...")
				activity.RecordHeartbeat(ctx, "")
			}
		}
	}()

	logger.Info("invoked activity")

	return activity.ErrResultPending
}
