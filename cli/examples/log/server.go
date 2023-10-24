package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/ubntc/go/cli/cli"
	"github.com/ubntc/go/cli/cli/config"
	"github.com/ubntc/go/cli/loggers/stdlogger"
)

// Server is an dummy server.
type Server struct {
	sync.WaitGroup
}

// Serve starts the server to test standard log usage.
func (s *Server) Serve(ctx context.Context) {
	s.Add(1)
	defer s.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("server stopped")
			return
		case <-time.After(1000 * time.Millisecond):
			log.Println("server alive")
		}
	}
}

// Shutdown waits for the server to stop.
func (s *Server) Shutdown() {
	s.Wait()
}

func main() {
	var (
		interactive = flag.Bool("i", false, "interactive mode")
		showClock   = flag.Bool("c", false, "display the live clock")
		verbose     = flag.Bool("v", false, "verbose logging")
	)
	flag.Parse()

	cfg := config.Default(*interactive)
	cfg.ShowClock = *showClock

	if *interactive {
		cli.SetupLogging(stdlogger.Setup)
	} else {
		log.SetFlags(log.LstdFlags | log.LUTC) // use unix time if running on non-interactive server
		log.Println("setting logger to UTC")
	}

	cli.GetTerm().SetVerbose(*verbose)

	ctx, cancel := cli.StartTerm(context.Background(), cfg)
	defer cancel()

	var srv Server
	go srv.Serve(ctx)

	<-ctx.Done()
	srv.Shutdown()
}
