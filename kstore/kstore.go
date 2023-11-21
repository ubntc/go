package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ubntc/go/kstore/examples"
	"github.com/ubntc/go/kstore/kstore"
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/provider/kafkago"
)

func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

const usageTitle = `
Tables is a Kafka topic manager that manages 'tables' in Kafka as compacted topics.

Usage:
`

func main() {
	f := flag.CommandLine
	f.Usage = func() {
		fmt.Fprint(f.Output(), usageTitle)
		flag.PrintDefaults()
	}
	var (
		reset   = flag.Bool("reset", false, "reset table schemas to the empty schema (runs first)")
		delete  = flag.Bool("delete", false, "delete all table topics (runs after reset)")
		purge   = flag.Bool("purge", false, "-reset and -delete together (runs first)")
		demo    = flag.Bool("demo", false, "manage some tables")
		cleanup = flag.Bool("cleanup", false, "-reset and -delete together (runs last)")
		timeout = flag.Duration("timeout", time.Second*10, "demo timeout")
		quiet   = flag.Bool("quiet", false, "more logging")
		group   = flag.String("group", "kstore", "group ID")
	)
	flag.Parse()

	cfg, err := config.LoadConfig()
	exitOnError(err)

	log.Println("bootstrap broker:", cfg.Server)

	groupConfig := config.Group{
		ID:     *group,
		Topics: []string{config.DefaultSchemasTopic},
	}

	c := kafkago.NewClient(cfg, config.DefaultProperties(), groupConfig)
	defer c.Close()
	tm := kstore.NewSchemaManager(config.DefaultSchemasTopic, c)
	if !*quiet {
		c.SetLogger(kafkago.NewLogger("SchemaManager"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = tm.Setup(ctx)
	exitOnError(err)

	if *purge {
		Purge(ctx, tm)
	} else {
		if *reset {
			Reset(ctx, tm)
		}
		if *delete {
			Delete(ctx, tm)
		}
	}

	if *demo {
		ctx, cancel := context.WithTimeout(ctx, *timeout)
		defer cancel()
		Demo(ctx, tm, cfg)
	}

	if *cleanup {
		Purge(ctx, tm)
	}
}

// Commands and Scenarios
// ======================

func Purge(ctx context.Context, tm *kstore.SchemaManager) {
	Reset(ctx, tm)
	Delete(ctx, tm)
}

func Reset(ctx context.Context, tm *kstore.SchemaManager) {
}

func Delete(ctx context.Context, tm *kstore.SchemaManager) {
}

func Demo(ctx context.Context, tm *kstore.SchemaManager, cfg *config.KeyFile) {
	exitOnError(examples.RunTopicManagement(ctx, tm, cfg))
}
