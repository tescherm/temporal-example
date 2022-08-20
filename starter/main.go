package main

import (
	"context"
	"github.com/pborman/uuid"
	"github.com/temporalio/samples-go/foo"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "pick-first_" + uuid.New(),
		TaskQueue: "pick-first",

		WorkflowTaskTimeout: 1 * time.Hour,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, example.ExampleWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	err = c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), "init", nil)
	if err != nil {
		panic(err)
	}
}
