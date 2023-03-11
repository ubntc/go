//go:generate weaver generate

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/ubntc/go/games/distris/client"
	"github.com/ubntc/go/games/distris/server"
)

const (
	SERVER  = "server"
	CLIENT  = "client"
	ADDRESS = "localhost:12345"
)

var modeHelp = fmt.Sprintf("[%s|%s]", SERVER, CLIENT)

func main() {
	var mode = flag.String("mode", CLIENT, modeHelp)
	var address = flag.String("address", ADDRESS, "server addesss")

	flag.Parse()
	var err error
	switch *mode {
	case SERVER:
		err = server.Run(*address)
	case CLIENT:
		err = client.Run(*address)
	default:
		err = errors.Errorf("invalid mode: %s, allowed values: %s", *mode, modeHelp)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
