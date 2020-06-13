package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/loggers/zerologger"
)

var interactive = flag.Bool("i", false, "interactive mode")
var debug = flag.Bool("debug", false, "debug mode")
var noClock = flag.Bool("n", false, "don't display clock")
var verbose = flag.Bool("v", false, "vebose output and prompts")
var demo = flag.String("demo", "", "script sequence of input keys and delays")

// Server is an dummy server.
type Server struct {
	sync.WaitGroup
	sync.RWMutex
	status      string
	logInterval time.Duration
}

// Serve starts the server to test zerolog usage.
func (s *Server) Serve(ctx context.Context) {
	s.Add(1)
	defer s.Done()
	for {
		s.SetStatus("idle")
		select {
		case <-ctx.Done():
			s.SetStatus("dead")
			s.Status()
			return
		case <-time.After(s.logInterval):
			s.SetStatus("active")
			s.Status()
			time.Sleep(s.logInterval / 5)
		}
	}
}

// Shutdown waits for the server to stop.
func (s *Server) Shutdown() {
	s.Wait()
}

// SetStatus sets the server status.
func (s *Server) SetStatus(status string) {
	s.Lock()
	defer s.Unlock()
	s.status = status
}

// Status logs the server status.
func (s *Server) Status() {
	s.RLock()
	defer s.RUnlock()
	log.Print("server is " + s.status)
}

func help() {
	t := cli.GetTerm()
	t.SetMessage(cli.GetCommands().String())
	t.Println(cli.GetCommands().Help())
}

func main() {
	flag.Parse()
	cli.GetTerm().SetDebug(*debug)
	cli.GetTerm().SetVerbose(*verbose)

	var opt []cli.Option
	srv := &Server{logInterval: 900 * time.Millisecond}
	if len(*demo) > 0 {
		srv.logInterval = time.Minute
	}

	if *noClock {
		opt = append(opt, cli.WithoutClock())
	}

	if *interactive {
		cli.SetupLogging(zerologger.Setup)
		opt = append(opt, cli.WithInput(cli.Commands{
			{Name: "help", Key: 'h', Fn: help},
			{Name: "status", Key: 's', Fn: srv.Status},
		}))
	} else {
		// use unix time if running on non-interactive server
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	ctx, cancel := cli.WithSigWait(context.Background(), opt...)
	defer cancel()

	go cli.GetCommands().RunScript(*demo)
	go srv.Serve(ctx)

	<-ctx.Done()
	srv.Shutdown()
}
