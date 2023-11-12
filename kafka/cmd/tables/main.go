package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ubntc/go/kafka/internal/cloudtables"
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
		verbose = flag.Bool("verbose", false, "more logging")
	)
	flag.Parse()

	cfg, err := cloudtables.LoadConfig()
	exitOnError(err)

	log.Println("bootstrap broker:", cfg.Server)

	tm := &cloudtables.TableManager{
		Topic:  cloudtables.DefaultManagerTopic,
		Writer: cloudtables.NewWriter(cfg),
	}

	if *verbose {
		tm.Writer.Logger = cloudtables.NewLogger("TableManager")
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
		Demo(ctx, tm)
	}

	if *cleanup {
		Purge(ctx, tm)
	}
}

// Commands and Scenarios
// ======================

func Purge(ctx context.Context, tm *cloudtables.TableManager) {
	Reset(ctx, tm)
	Delete(ctx, tm)
}

func Reset(ctx context.Context, tm *cloudtables.TableManager) {
	name := "table1"
	_ = tm.CreateOrUpdateTable(ctx, cloudtables.Table{Name: name})
}

func Delete(ctx context.Context, tm *cloudtables.TableManager) {
	name := "table1"
	_ = tm.DeleteTable(ctx, name)
}

func Demo(ctx context.Context, tm *cloudtables.TableManager) {
	tbl := cloudtables.Table{
		Name: "table1",
		Schema: []cloudtables.Field{
			{Name: "col1", Type: cloudtables.FieldTypeString},
		},
	}

	exitOnError(tm.CreateOrUpdateTable(ctx, tbl))

	tbl.Schema = append(tbl.Schema, cloudtables.Field{Name: "col2", Type: cloudtables.FieldTypeString})
	exitOnError(tm.CreateOrUpdateTable(ctx, tbl))

	tbl.Schema = append(tbl.Schema, cloudtables.Field{Name: "col3", Type: cloudtables.FieldTypeString})
	exitOnError(tm.CreateOrUpdateTable(ctx, tbl))

	cfg, err := cloudtables.LoadConfig()
	exitOnError(err)

	w := cloudtables.NewWriter(cfg)
	w.Topic = tm.TopicForTable(tbl.Name)

	err = w.WriteMessages(ctx, cloudtables.GenerateMessages(tbl, 10)...)
	exitOnError(err)

	log.Println("wrote 10 messages to table:", tbl.Name)
}
