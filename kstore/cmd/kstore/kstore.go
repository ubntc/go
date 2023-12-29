package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ubntc/go/kstore/examples"
	"github.com/ubntc/go/kstore/kstore/cli"
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/kstore/manager"
	"github.com/ubntc/go/kstore/provider/api"
	"github.com/ubntc/go/kstore/provider/kafkago"
)

func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getClient(cfg *config.KeyFile, group config.Group) api.Client {
	return kafkago.NewClient(cfg, config.DefaultProperties(), group)
}

func main() {
	// 1. setup custom flags and commands
	var (
		// commands to run on the cluster
		cleanup = flag.Bool("cleanup", false, "complete the run with reset and delete even after error")

		// chain of commands to run
		timeout = flag.Duration("timeout", time.Second*10, "demo timeout")

		demo = manager.Action{
			Name: "demo",
			Help: "run the demo",
			Func: manager.WithTimeout(examples.RunTopicManagement, *timeout),
		}
	)

	// 2. hand control over to the kstore CLI and let it parse and setup all commands
	workflow, err := cli.Parse(getClient, demo)
	exitOnError(err)

	// apply custom flags
	if *cleanup {
		workflow.DeferStep(manager.Purge)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exitOnError(workflow.Run(ctx, manager.OnErrorStop))
}
