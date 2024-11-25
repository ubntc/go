package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/cli/config"
	"github.com/ubntc/go/cli/loggers/zerologger"
)

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
			s.LogStatus(ctx)
			return
		case <-time.After(s.logInterval):
			s.SetStatus("active")
			s.LogStatus(ctx)
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

// LogStatus logs the server status.
func (s *Server) LogStatus(context.Context) {
	s.RLock()
	defer s.RUnlock()
	log.Print("server is " + s.status)
}

func help(context.Context) {
	t := cli.GetTerm()
	t.SetMessage(cli.GetCommands().String())
	t.Println(cli.GetCommands().Help())
}

func main() {
	var (
		interactive = flag.Bool("i", false, "interactive mode")
		debug       = flag.Bool("debug", false, "debug mode")
		showClock   = flag.Bool("c", false, "show the live clock")
		verbose     = flag.Bool("v", false, "vebose output and prompts")
		demo        = flag.String("demo", "", "script sequence of input keys and delays")
	)

	flag.Parse()
	cli.GetTerm().SetDebug(*debug)
	cli.GetTerm().SetVerbose(*verbose)

	srv := &Server{logInterval: 900 * time.Millisecond}
	if len(*demo) > 0 {
		srv.logInterval = time.Minute
	}

	var cmds cli.Commands
	cfg := config.Default(*interactive)
	cfg.ShowClock = *showClock

	if *interactive {
		cli.SetupLogging(zerologger.Setup)
		cmds = cli.Commands{
			{Name: "help", Key: 'h', Fn: help},
			{Name: "status", Key: 's', Fn: srv.LogStatus},
		}
	} else {
		// use unix time if running on non-interactive server
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	ctx, cancel := cli.StartTerm(context.Background(), cfg, cmds...)
	defer cancel()

	go cli.GetCommands().RunScript(*demo) // nolint
	go srv.Serve(ctx)

	<-ctx.Done()
	srv.Shutdown()
}
