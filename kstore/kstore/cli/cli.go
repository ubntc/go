package cli

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ubntc/go/kstore/kstore"
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/kstore/manager"
	"github.com/ubntc/go/kstore/provider/api"
)

var shortInfo = "KStore manages evolvable data tables as compacted topics in Kafka."

var longInfo = strings.TrimSpace(`
==========
KStore CLI
==========
` + shortInfo + `
For more see: github.com/ubntc/go/kstore.

The KStore CLI is used to create, update, and delete these tables and their schemas.
It can be customized to run read/write scenarios as part of an application and
is used for KStore's own comprehensive test scenarios.
`)

var usageTemplate = longInfo + `
Usage: kstore [FLAGS] [COMMANDS...]

Commands: %s

%s

Flags:

`

func usageLines(commands []manager.Action) (names, desc []string) {
	maxLen := 0
	tab := "  "
	for _, v := range commands {
		names = append(names, v.Name)
		maxLen = max(maxLen, len(v.Name))
	}
	for _, v := range commands {
		spacer := strings.Repeat(" ", maxLen-len(v.Name))
		line := tab + v.Name + spacer + tab + v.Help
		desc = append(desc, line)
	}
	return
}

// ClientGetter creates a new client.
// It is called by the CLI after the secrets have been successfully loaded to setup the SchemaManager.
type ClientGetter func(cfg *config.KeyFile, group config.Group) api.Client

// Parse parses all CLI arguments and uses them to setup a complete Workflow.
func Parse(getClient ClientGetter, customActions ...manager.Action) (*manager.Workflow, error) {
	actions := append(manager.Actions(), customActions...)

	f := flag.CommandLine
	f.Usage = func() {
		names, desc := usageLines(actions)
		usage := fmt.Sprintf(usageTemplate, strings.Join(names, "|"), strings.Join(desc, "\n"))
		fmt.Fprint(f.Output(), usage)
		flag.PrintDefaults()
	}
	var (
		// chain of commands to run
		verbose    = flag.Bool("verbose", false, "more logging")
		group      = flag.String("group", "kstore", "group ID")
		groupShort = flag.String("g", "kstore", "group ID (short form of -group)")
		table      = flag.String("table", "", "ID of the managed table")
		tableShort = flag.String("t", "", "ID of the managed table (short form of -table)")
		all        = flag.Bool("all", false, "must be set to run an operation on KStore-managed ALL tables in the cluster")
	)
	flag.Parse()

	// input validation
	if *table == "" && *tableShort != "" {
		*table = *tableShort
	}
	if *group == "" && *groupShort != "" {
		*group = *groupShort
	}
	if *all && *table != "" {
		return nil, fmt.Errorf("either -all XOR -t|-table must be set")
	}

	cfg, err := config.LoadKeyFile()
	if err != nil {
		return nil, err
	}

	log.Println("bootstrap broker:", cfg.Server)

	groupConfig := config.Group{
		ID:     *group,
		Topics: []string{config.DefaultSchemasTopic},
	}

	client := getClient(cfg, groupConfig)
	tm := manager.NewSchemaManager(config.DefaultSchemasTopic, client)

	if *verbose {
		// TODO: extract logger interface
		client.SetLogger(kstore.NewLogger("SchemaManager"))
	}

	program := flag.Args()

	// setup workflow
	return manager.NewWorkflow(tm, cfg, actions, program)
}
